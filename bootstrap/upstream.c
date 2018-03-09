#include <syslog.h>

extern void switchuser();
extern void daemonise();
extern void coreprocess(char **argv);

int main(int argc, char **argv)
{
    openlog("achelous/upstream", LOG_PID, LOG_MAIL);

    daemonise();
    switchuser();
    coreprocess(argv);

    return 0;
}
