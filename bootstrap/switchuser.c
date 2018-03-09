#include <unistd.h>
#include <pwd.h>
#include <grp.h>
#include <syslog.h>
#include <errno.h>

#include "config.h"

void switchuser()
{
    struct group* g = getgrnam(CONFIG_GROUP);
    if(g == NULL)
    {
        syslog(LOG_ERR, "failed to find group");
        _exit(1);
    }

    struct passwd* u = getpwnam(CONFIG_USER);
    if(u == NULL)
    {
        syslog(LOG_ERR, "failed to find user");
        _exit(1);
    }

    if(setgid(g->gr_gid) < 0)
    {
        syslog(LOG_ERR, "failed to set gid (errno=%d)", errno);
        _exit(1);
    }

    if(setuid(u->pw_uid) < 0)
    {
        syslog(LOG_ERR, "failed to set uid (errno=%d)", errno);
        _exit(1);
    }
}
