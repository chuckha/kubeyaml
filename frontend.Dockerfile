FROM nginx:mainline
ADD dist /srv/www
ADD infra/kubeyaml.com.conf /etc/nginx/conf.d/default.conf

