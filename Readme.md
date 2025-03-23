# OU/User 追加
```
$ docker exec -it openldap /bin/bash

$ ldapadd -x -D "cn=admin,dc=example,dc=com" -w adminpassword <<EOF
dn: ou=devops,dc=example,dc=com
objectClass: organizationalUnit
ou: devops
EOF

$ ldapadd -x -D "cn=admin,dc=example,dc=com" -w adminpassword <<EOF
dn: cn=amrutha,ou=devops,dc=example,dc=com
objectClass: inetOrgPerson
cn: amrutha
sn: Amrutha
uid: amrutha
userPassword: Amrutha@123
EOF

```
