```sh
echo "Lorem body..." | mail -r "sender@mail.com" -s "Lorem subject..." "receiver@mail.com"
echo "Lorem body..." | strace -f -e trace=process mail -r "sender@mail.com" -s "Lorem subject..." "receiver@mail.com"
```
