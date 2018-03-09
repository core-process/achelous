#include <stdio.h>

int readpid (char *pidfile)
{
    FILE *f;
    int pid = 0;

    if (!(f=fopen(pidfile,"r")))
        return 0;
    fscanf(f,"%d", &pid);
    fclose(f);
    return pid;
}
