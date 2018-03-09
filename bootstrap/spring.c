#include <syslog.h>

extern void switchuser();
extern void coreprocess(char **argv);

int main(int argc, char **argv)
{
    // setup logging
    openlog("achelous/spring", LOG_PERROR | LOG_PID, LOG_MAIL);

    // switch to mailing user
    switchuser();

    // execute core process (does not return)
    coreprocess(argv);
    return 0;
}
