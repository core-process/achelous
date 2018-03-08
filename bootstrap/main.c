#include <sys/types.h>
#include <pwd.h>
#include <grp.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <linux/limits.h>

#include "config.h"

int switch_user()
{
    struct group* g = getgrnam(CONFIG_GROUP);
    if(g == NULL)
    {
        errno = EINVAL;
        return -1;
    }

    struct passwd* u = getpwnam(CONFIG_USER);
    if(u == NULL)
    {
        errno = EINVAL;
        return -1;
    }

    if(setgid(g->gr_gid) < 0)
    {
        return -1;
    }

    if(setuid(u->pw_uid) < 0)
    {
        return -1;
    }

    return 0;
}

int run_core(char** argv)
{
    // read path of current process
    char path[PATH_MAX];
    memset(path, 0, sizeof(path));

    if (readlink("/proc/self/exe", path, sizeof(path)-1) < 0)
    {
        return -1;
    }

    // add core extension
    strncat(path, "-core", PATH_MAX -strlen(path) -1);

    // run core process
    if(execv(path, argv) < 0)
    {
        return -1;
    }

    return 0;
}

int main(int argc, char** argv)
{
    if(switch_user() < 0)
    {
        perror("switch_user");
        return -1;
    }

    if(run_core(argv) < 0)
    {
        perror("run_core");
        return -1;
    }

    return 0;
}
