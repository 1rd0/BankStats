# SRETest
## Условия 
1) 3 зоны/5 нод,
2) 5–10 c старт,
3) 4 пода держат пик,
4) CPU бурст, 128Mi,
5) дневной цикл


## Запуск всех манифестов

Находясь в папке с файлами пропишем:
 ```bash
 
kubectl apply -f . 
 ```
Так же для проверки 
```bash 

kubectl get all
 ```
  


## Что сделано

1) Deployment: 4 пода, с распределением по зонам (через topologySpreadConstraints).

2) Probes(Код с пробированием закоментирован из-за отсутсвия реализации поинтов).

3) Resources: Requests 0.1 CPU / 128Mi, Limits до 1 CPU / 256Mi для пиковых нагрузок.

4) SecurityContext: без root, без эскалации прав.

5) Lifecycle: preStop для graceful shutdown.

6) PodDisruptionBudget: минимум 50% подов всегда доступны.

7) HPA: авто-скейл подов по CPU.
   

<p align="center">
 <img   src="img.png" alt="qr"/>
</p>
