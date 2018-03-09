#include <syslog.h>
#include <stddef.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <bsd/libutil.h>

#include "config.h"

extern void switchuser();
extern void daemonise();
extern void coreprocess(char **argv);

void cmdstart(int daemon, char **argv)
{
    // create lock file
    struct pidfh *pfh = pidfile_open(CONFIG_UPID, 0600, NULL);
    if (pfh == NULL)
    {
        syslog(LOG_ERR, "failed to lock pidfile (errno=%d)", errno);
        _exit(1);
    }

    // become a daemon
    if(daemon != 0)
    {
        daemonise();
    }

    // write pid
    if(pidfile_write(pfh) == -1)
    {
        syslog(LOG_ERR, "failed to write pid to pidfile (errno=%d)", errno);
        _exit(1);
    }

    // create child for core process
    pid_t pid = fork();

    if (pid == -1)
    {
        syslog(LOG_ERR, "failed to fork for core process (errno=%d)", errno);
        _exit(1);
    }

    if(pid == 0)
    {
        // close pid file in core process
        if(pidfile_close(pfh) == -1)
        {
            syslog(LOG_ERR, "failed to close pidfile in core process (errno=%d)", errno);
            _exit(1);
        }

        // switch to mailing user
        syslog(LOG_INFO, "switching to mailing user");
        switchuser();

        // execute core process (does not return)
        syslog(LOG_INFO, "starting core process");
        coreprocess(argv);
    }

    // wait for core process to end
    if(waitpid(pid, NULL, 0) == -1)
    {
        syslog(LOG_ERR, "failed to wait for core process (errno=%d)", errno);
        _exit(1);
    }

    // remove pid file
    if(pidfile_remove(pfh) == -1)
    {
        syslog(LOG_ERR, "failed to remove pidfile (errno=%d)", errno);
        _exit(1);
    }

    syslog(LOG_INFO, "completed successfuly");
    _exit(0);
}

void cmdstop()
{
}

int main(int argc, char **argv)
{
    if(argc < 2)
    {
        // setup logging
        openlog("achelous/upstream", LOG_PERROR | LOG_PID, LOG_MAIL);

        // run in foreground mode (does not return)
        cmdstart(0, argv);
    }
    else
    {
        // setup logging
        openlog("achelous/upstream", LOG_PID, LOG_MAIL);

        // start daemon (does not return)
        if(strcmp(argv[1], "start") == 0)
        {
            cmdstart(1, argv);
        }

        // stop daemon (does not return)
        if(strcmp(argv[1], "stop") == 0)
        {
            cmdstop();
        }
    }

    return 0;
}
