#include <unistd.h>
#include <pwd.h>
#include <grp.h>
#include <syslog.h>
#include <errno.h>

#include "config.h"

void switchuser()
{
    struct group *g = getgrnam(CONFIG_GROUP);
    if (g == NULL)
    {
        syslog(LOG_ERR, "failed to find group %s", CONFIG_GROUP);
        _exit(1);
    }

    struct passwd *u = getpwnam(CONFIG_USER);
    if (u == NULL)
    {
        syslog(LOG_ERR, "failed to find user %s", CONFIG_USER);
        _exit(1);
    }

    if (setgid(g->gr_gid) < 0)
    {
        syslog(LOG_ERR, "failed to set gid %d (errno=%d)", g->gr_gid, errno);
        _exit(1);
    }

    if (setuid(u->pw_uid) < 0)
    {
        syslog(LOG_ERR, "failed to set uid %d (errno=%d)", u->pw_uid, errno);
        _exit(1);
    }
}
