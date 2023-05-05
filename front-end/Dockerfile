# Imagem base do Ubuntu
FROM ubuntu:latest

# Atualiza o sistema e instala o Apache
RUN apt-get update && \
    apt-get install -y apache2 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copia os arquivos do diretório raiz para o diretório padrão do Apache
COPY . /var/www/html/

# Expõe a porta padrão do Apache (80)
EXPOSE 80

# Inicia o Apache quando o container for executado
CMD ["/usr/sbin/apache2ctl", "-D", "FOREGROUND"]