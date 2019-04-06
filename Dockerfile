FROM nginx:mainline
ADD dist /www/data
ADD infra/kubeyaml.com.conf /etc/nginx/conf.d/default.conf