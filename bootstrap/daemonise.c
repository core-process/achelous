#include <unistd.h>
#include <syslog.h>
#include <signal.h>
#include <errno.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <time.h>
#include <bsd/libutil.h>

extern int readpid(const char *pidfile);

struct pidfh* daemonise(const char* pidpath)
{
    // lock pid file
    struct pidfh* pfh = pidfile_open(pidpath, 0600, NULL);
    if (pfh == NULL)
    {
        syslog(LOG_ERR, "failed to lock pid file (errno=%d)", errno);
        goto main_exit_fail;
    }

    // The first call to fork(2) ensures that the process is
    // not a group leader, so that it is possible for that
    // process to create a new session and become a session
    // leader. There are other reasons for the first
    // fork(2): if the daemon was started as a shell
    // command, having the process fork and the parent exit
    // makes the shell go back to its prompt and wait for
    // more commands.
    pid_t pid = fork();
    if (pid == -1)
    {
        syslog(LOG_ERR, "failed to fork while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }
    else if (pid != 0)
    {
        // wait 10 sec for pid file to be written
        for(time_t ts = time(NULL); (time(NULL) - ts) <= 10;)
        {
            sleep(1);

            int observed_pid = readpid(pidpath);
            if(observed_pid != 0)
            {
                syslog(LOG_ERR, "observed pid: %d", observed_pid);
                goto branch_exit_grace;
            }
        }

        // fail, since we did not observe pid to be written
        syslog(LOG_ERR, "failed to observe pid");
        goto branch_exit_fail;
    }

    // Start a new session for the daemon.
    if (setsid() == -1)
    {
        syslog(LOG_ERR, "failed to become a session leader while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }

    // The second fork(2) is there to ensure that the new
    // process is not a session leader, so it won't be able
    // to (accidentally) allocate a controlling terminal,
    // since daemons are not supposed to ever have a
    // controlling terminal.
    signal(SIGHUP, SIG_IGN);

    pid = fork();
    if (pid == -1)
    {
        syslog(LOG_ERR, "failed to fork while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }
    else if (pid != 0)
    {
        goto branch_exit_grace;
    }

    // Set the current working directory to the root directory.
    if (chdir("/") == -1)
    {
        syslog(LOG_ERR, "failed to change working directory while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }

    // Set the user file creation mask to zero.
    umask(0);

    // Close then reopen standard file descriptors.
    int fd;

    if ((fd = open("/dev/null", O_RDONLY)) == -1)
    {
        syslog(LOG_ERR, "failed to reopen stdin while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }

    dup2(fd, 0);
    close(fd);

    if ((fd = open("/dev/null", O_WRONLY)) == -1)
    {
        syslog(LOG_ERR, "failed to reopen stdout while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }

    dup2(fd, 1);
    close(fd);

    if ((fd = open("/dev/null", O_WRONLY)) == -1)
    {
        syslog(LOG_ERR, "failed to reopen stderr while daemonising (errno=%d)", errno);
        goto main_exit_fail;
    }

    dup2(fd, 2);
    close(fd);

    // return pid file handle
    return pfh;

    ////////////////////////////////////////////////////////
    // EXIT HANDLING
    ////////////////////////////////////////////////////////

main_exit_fail:
    if(pfh != NULL)
    {
        pidfile_remove(pfh);
        syslog(LOG_INFO, "pid file removed");
    }
    _exit(1);

branch_exit_fail:
    if(pfh != NULL)
    {
        pidfile_close(pfh);
    }
    _exit(1);

branch_exit_grace:
    if(pfh != NULL)
    {
        pidfile_close(pfh);
    }
    _exit(0);
}
