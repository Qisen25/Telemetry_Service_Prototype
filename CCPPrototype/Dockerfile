FROM node:latest

COPY . /src

WORKDIR /src

RUN npm install

EXPOSE 3000

CMD [ "node", "ReceivingWS.js"]
