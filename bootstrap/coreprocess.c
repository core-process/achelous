#include <unistd.h>
#include <syslog.h>
#include <errno.h>
#include <string.h>
#include <linux/limits.h>

void coreprocess(char** argv)
{
    // read path of current process
    char path[PATH_MAX];
    memset(path, 0, sizeof(path));

    if (readlink("/proc/self/exe", path, sizeof(path)-1) < 0)
    {
        syslog(LOG_ERR, "failed to read executable path (errno=%d)", errno);
        _exit(1);
    }

    // add core extension
    strncat(path, "-core", PATH_MAX -strlen(path) -1);

    // run core process
    if(execv(path, argv) < 0)
    {
        syslog(LOG_ERR, "failed to execute core process (errno=%d)", errno);
        _exit(1);
    }
}
