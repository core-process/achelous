#include <syslog.h>

extern void switchuser();
extern void daemonise();
extern void coreprocess(char** argv);

int main(int argc, char** argv)
{
    openlog("achelous-upstream", LOG_PERROR|LOG_PID, LOG_MAIL);

    switchuser();
    daemonise();
    coreprocess(argv);

    return 0;
}
