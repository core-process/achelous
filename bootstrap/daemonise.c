#include <unistd.h>
#include <syslog.h>
#include <signal.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

void daemonise()
{
    // Fork, allowing the parent process to terminate.
    pid_t pid = fork();
    if (pid == -1)
    {
        syslog(LOG_ERR, "failed to fork while daemonising (errno=%d)", errno);
        _exit(1);
    }
    else if (pid != 0)
    {
        _exit(0);
    }

    // Start a new session for the daemon.
    if (setsid() == -1)
    {
        syslog(LOG_ERR, "failed to become a session leader while daemonising (errno=%d)", errno);
        _exit(1);
    }

    // Fork again, allowing the parent process to terminate.
    signal(SIGHUP, SIG_IGN);
    pid = fork();
    if (pid == -1)
    {
        syslog(LOG_ERR, "failed to fork while daemonising (errno=%d)", errno);
        _exit(1);
    }
    else if (pid != 0)
    {
        _exit(0);
    }

    // Set the current working directory to the root directory.
    if (chdir("/") == -1)
    {
        syslog(LOG_ERR, "failed to change working directory while daemonising (errno=%d)", errno);
        _exit(1);
    }

    // Set the user file creation mask to zero.
    umask(0);

    // Close then reopen standard file descriptors.
    int fd;

    if ((fd = open("/dev/null", O_RDONLY)) == -1)
    {
        syslog(LOG_ERR, "failed to reopen stdin while daemonising (errno=%d)", errno);
        _exit(1);
    }
    dup2(fd, 0);
    close(fd);

    if ((fd = open("/dev/null", O_WRONLY)) == -1)
    {
        syslog(LOG_ERR, "failed to reopen stdout while daemonising (errno=%d)", errno);
        _exit(1);
    }
    dup2(fd, 1);
    close(fd);

    if ((fd = open("/dev/null", O_WRONLY)) == -1)
    {
        syslog(LOG_ERR, "failed to reopen stderr while daemonising (errno=%d)", errno);
        _exit(1);
    }
    dup2(fd, 2);
    close(fd);
}
