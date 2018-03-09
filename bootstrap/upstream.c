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
extern int readpid (char *pidfile);

pid_t corepid = -1;

void signal2core(int signum)
{
    syslog(LOG_INFO, "forwarding signal to core (signum=%d)", signum);
    kill(corepid, signum);
}

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
    corepid = fork();

    if (corepid == -1)
    {
        syslog(LOG_ERR, "failed to fork for core process (errno=%d)", errno);
        _exit(1);
    }

    if(corepid == 0)
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
    signal(SIGINT, signal2core);

    if(waitpid(corepid, NULL, 0) == -1)
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
    // read pid
    int pid = readpid(CONFIG_UPID);

    if(!pid)
    {
        syslog(LOG_ERR, "failed to read pid from file");
        _exit(1);
    }

    // kill process
    if(kill(pid, SIGINT) == -1)
    {
        syslog(LOG_ERR, "failed to kill pid %d (errno=%d)", pid, errno);
        _exit(1);
    }

    syslog(LOG_INFO, "send kill signal successfuly");
    _exit(0);
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
