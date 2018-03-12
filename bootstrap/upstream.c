#include <syslog.h>
#include <stddef.h>
#include <errno.h>
#include <string.h>
#include <unistd.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <bsd/libutil.h>

#include "config.h"

extern void switchuser();
extern struct pidfh* daemonise(const char* pidpath);
extern void coreprocess(char **argv);
extern int readpid(const char *pidfile);

pid_t corepid = -1;

void SIG_CORE(int signum)
{
    syslog(LOG_INFO, "propagating signal %d to core process", signum);
    if(corepid != -1)
    {
        kill(corepid, signum);
    }
}

void cmdstart(char **argv)
{
    // become a daemon (returns on daemonized process)
    struct pidfh *pfh = daemonise(CONFIG_UPSTREAM_PID);

    // fork controlling and core process
    pid_t pid = fork();
    if (pid == -1)
    {
        syslog(LOG_ERR, "failed to fork for core process (errno=%d)", errno);
        goto exit_fail;
    }
    else if (pid != 0)
    {
        // remember pid of core process
        corepid = pid;

        // write pid
        if (pidfile_write(pfh) == -1)
        {
            syslog(LOG_ERR, "failed to write pid to pid file (errno=%d)", errno);
            goto exit_fail;
        }
        syslog(LOG_INFO, "pid file written");

        // wait for core process to end
        signal(SIGINT, SIG_CORE);
        signal(SIGTERM, SIG_CORE);
        signal(SIGQUIT, SIG_CORE);
        signal(SIGUSR1, SIG_CORE);

        if (waitpid(corepid, NULL, 0) == -1)
        {
            syslog(LOG_ERR, "failed to wait for core process (errno=%d)", errno);
            goto exit_fail;
        }

        // service completed
        syslog(LOG_INFO, "service completed");
        goto exit_grace;
    }

    // close pid file
    pidfile_close(pfh);

    // switch to mailing user
    switchuser();

    // execute core process (does not return)
    syslog(LOG_INFO, "executing core process");
    coreprocess(argv);
    return;

    ////////////////////////////////////////////////////////
    // EXIT HANDLING (controlling process)
    ////////////////////////////////////////////////////////

exit_fail:
    if(pfh != NULL)
    {
        pidfile_remove(pfh);
        syslog(LOG_INFO, "pid file removed");
    }
    _exit(1);

exit_grace:
    if(pfh != NULL)
    {
        pidfile_remove(pfh);
        syslog(LOG_INFO, "pid file removed");
    }
    _exit(0);
}

void cmdstop()
{
    // read pid
    int pid = readpid(CONFIG_UPSTREAM_PID);

    if (!pid)
    {
        syslog(LOG_ERR, "failed to read pid from file");
        _exit(1);
    }

    // kill process
    if (kill(pid, SIGTERM) == -1)
    {
        syslog(LOG_ERR, "failed to send SIGTERM to pid %d (errno=%d)", pid, errno);
        _exit(1);
    }

    syslog(LOG_INFO, "sent SIGTERM to service");
    _exit(0);
}

void cmdstatus()
{
    // read pid
    int pid = readpid(CONFIG_UPSTREAM_PID);

    // print status
    if (!pid)
    {
        syslog(LOG_INFO, "service is stopped");
        printf("stopped\n");
    }
    else
    {
        syslog(LOG_INFO, "service is running (pid=%d)", pid);
        printf("running [%d]\n", pid);
    }
    fflush(stdout);
    _exit(0);
}

void cmdtrigger()
{
    // read pid
    int pid = readpid(CONFIG_UPSTREAM_PID);

    if (!pid)
    {
        syslog(LOG_ERR, "failed to read pid from file");
        _exit(1);
    }

    // send SIGUSR1 to process
    if (kill(pid, SIGUSR1) == -1)
    {
        syslog(LOG_ERR, "failed to send SIGUSR1 to pid %d (errno=%d)", pid, errno);
        _exit(1);
    }

    syslog(LOG_INFO, "sent SIGUSR1 to service");
    _exit(0);
}

int main(int argc, char **argv)
{
    // setup logging
    openlog("achelous/upstream", LOG_PID, LOG_MAIL);

    // start daemon (does not return)
    if (argc < 2 || strcmp(argv[1], "start") == 0)
    {
        cmdstart(argv);
    }

    // stop daemon (does not return)
    if (strcmp(argv[1], "stop") == 0)
    {
        cmdstop();
    }

    // daemon status (does not return)
    if (strcmp(argv[1], "status") == 0)
    {
        cmdstatus();
    }

    // trigger queue run (does not return)
    if (strcmp(argv[1], "trigger") == 0)
    {
        cmdtrigger();
    }

    return 0;
}
