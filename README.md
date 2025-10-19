# CONFIGURANDO O REDIS
crie um arquivo chamado redis-users.acl onde vai ter as credenciais de acesso ao redis com o conteúdo abaixo:

Obs: as vaŕiasvies são de exemplos, em produção use variávies mais robustas.

```acl
user default off
user admin on >minha_senha_admin allcommands allkeys
user app_user on >minha_senha_app +get +set +del +ping ~*
```

Depois disso converta para LF usando o dos2unix
```
sudo dos2unix redis-users.acl
```



# Rodando o docker compose