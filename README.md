# e-Financeira - SPED - Receita Federal

API para integração com módulo da e-Financeira do SPED da Receita Federal.

## Variáveis

.env
```shell
DB_URL=mongodb://mongo:27017
DB_NAME=efinanceira
DB_USERNAME=root
DB_PASSWORD=root

#SMTP
SMTP_HOST=
SMTP_PORT=465
SMTP_USERNAME=contato@overall.cloud
SMTP_PASSWORD=
```

Importante definir variaveis de ambiente com console. Exemplo:

```shell
export DB_URL="mongodb://mongo:27017"
```

## Ambiente de Desenvolvimento

[Quick Start do GoReleaser.](https://goreleaser.com/quick-start/)

### Configurando Docker

Pré configuração:

```shell
go mod init sped-efinanceira
go mod tidy
```

Criando a imagem:

```shell
docker-compose up
```

Para verificar (listagem containers):

```shell
docker ps
```

Para remover:

```shell
docker-compose down
```

Para logs:

```shell
docker logs -f backend-efinanceira mongo mongo-express
```

## Ambiente de Produção
    
 ### Instalanndo e Configurando no Servidor

Instalar o Go no Ubuntu:

 ```shell
sudo apt install golang-go git
 ```

### Clone o repositório

```shell
git clone git@github.com:gustavokennedy/sped-efinanceira.git && cd sped-efinanceira
```

### Instalando Dependências

```shell
go build main.go
```

<a href="https://www.digitalocean.com/community/tutorials/how-to-install-nginx-on-ubuntu-22-04" target="_Blank">Instalar o Nginx no Ubuntu.</a>

<a href="https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-22-04" target="_Blank">Instalar SSL com Nginx no Ubuntu.</a>

Primeiro, crie um novo arquivo no /lib/systemd/system/ chamado efinanceira.service:

 ```shell
 sudo nano /lib/systemd/system/efinanceira.service
 ```
 
 ```shell
[Unit]
Description=efinanceira

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/home/ubuntu/go/sped-efinanceira/main

[Install]
WantedBy=multi-user.target
```

Agora que você escreveu o arquivo da unidade de serviço, inicie seu serviço da web Go executando:

```shell
 sudo service efinanceira start
 ```

Para confirmar se o serviço está em execução, use o seguinte comando:

```shell
 sudo service efinanceira status
 ```

Para verificar no Log no Service:

  ```shell
 sudo journalctl -u efinanceira -b
 ```

 Para reiniciar configurações de Service:

  ```shell
 sudo systemctl daemon-reload
 ```

 ### Configurando Nginx

 Primeiro, altere seu diretório de trabalho para o sites-enabled do Nginx:

```shell
sudo nano /etc/nginx/sites-enabled/default
 ```

Adicione as seguintes linhas ao arquivo para estabelecer as configurações:

```shell
server {
    server_name _;

    location / {
        proxy_pass http://localhost:8080;
    }
}
 ```

Em seguida, recarregue suas configurações do Nginx executando o comando reload:

```shell
sudo nginx -s reload
 ```
