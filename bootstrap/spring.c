extern void switchuser();
extern void coreprocess(char** argv);

int main(int argc, char** argv)
{
    switchuser();
    coreprocess(argv);
    return 0;
}
