# desafio-multithreading-api
Este projeto implementa um cliente em Golang para busca de endereços a partir de um CEP, utilizando duas APIs públicas em chamadas concorrentes:  BrasilAPI  ViaCEP  A aplicação dispara requisições simultâneas para ambas as APIs e retorna apenas a resposta mais rápida, descartando a mais lenta.
