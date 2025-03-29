# OU/User 追加
```
$ docker exec -it openldap /bin/bash

$ ldapadd -x -D "cn=admin,dc=example,dc=com" -w adminpassword <<EOF
dn: ou=devops,dc=example,dc=com
objectClass: organizationalUnit
ou: devops
EOF

$ ldapadd -x -D "cn=admin,dc=example,dc=com" -w adminpassword <<EOF
dn: cn=ryotaro,ou=devops,dc=example,dc=com
objectClass: inetOrgPerson
cn: ryotaro
sn: Ryotaro
uid: ryotaro
userPassword: Ryotaro@123
EOF

```
