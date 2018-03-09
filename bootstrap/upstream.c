extern void switchuser();
extern void daemonise();
extern void coreprocess(char** argv);

int main(int argc, char** argv)
{
    switchuser();
    daemonise();
    coreprocess(argv);
    return 0;
}
