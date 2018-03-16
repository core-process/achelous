```sh
# EXAMPLE 1
sendmail -i "receiver@mail.com" <<EOF
Subject: Lorem subject...
From: sender@mail.com

Lorem body...
EOF

# EXAMPLE 2
# requires the GNU mail command (apt-get install mailutils)
echo "Lorem body..." | mail -r "sender@mail.com" -s "Lorem subject..." "receiver@mail.com"

# EXAMPLE 3
# using strace to detect call parameters to spring/sendmail
# (for debugging purposes)
echo "Lorem body..." | strace -f -e trace=process mail -r "sender@mail.com" -s "Lorem subject..." "receiver@mail.com"
```
