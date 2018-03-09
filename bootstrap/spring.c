#include <syslog.h>

extern void switchuser();
extern void coreprocess(char **argv);

int main(int argc, char **argv)
{
    openlog("achelous/spring", LOG_PERROR | LOG_PID, LOG_MAIL);

    switchuser();
    coreprocess(argv);

    return 0;
}
