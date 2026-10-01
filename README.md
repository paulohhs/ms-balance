# MS-Balance

Microsserviço **Balances** (consumidor Kafka) do desafio *EDA: Microsserviço de Balances* do curso FullCycle.
Ele consome os eventos `BalanceUpdated` publicados pelo **Wallet Core** (`wallet-core/`, produtor)
no tópico `balances`, mantém o saldo de cada conta em um MySQL próprio e expõe a consulta
`GET /balances/{account_id}` na porta **3003**.

## Arquitetura

```
POST /transactions ─► wallet-app ─► mysql-wallet
                          │
                          └─(Kafka: tópico "balances", evento BalanceUpdated)─► balance-app ─► mysql-balance
                                                                                   │
                                                          GET /balances/{id} ◄─────┘
```

## Pré-requisitos

- Docker e Docker Compose v2 (`docker compose`)
- Portas livres: 8080, 3003, 3306, 3307, 9092, 2181, 9021

## Como executar

```bash
docker compose up -d
```

Só isso. Na primeira subida (1–3 minutos), o compose:

1. sobe Zookeeper e Kafka e cria os tópicos `transactions` e `balances` (serviço `kafka-init`);
2. sobe os dois MySQL e roda **migrations e seeds** automaticamente (`wallet-core/sql/` e `sql/`);
3. sobe `wallet-app` (porta 8080) e `balance-app` (porta 3003) quando bancos e tópicos estão prontos.

Para acompanhar: `docker compose ps` e `docker compose logs -f balance-app wallet-app`.

## Testando

Abra `api/api.http` no VS Code (extensão **REST Client**) e execute em ordem:

1. `GET /balances/...` → saldos iniciais (John 1000, Jane 1000);
2. `POST /transactions` no wallet → transfere 100 do John para a Jane;
3. `GET /balances/...` → John 900, Jane 1100 (atualizado via Kafka).

### Dados de seed

| Cliente | account_id | Saldo |
|---|---|---|
| John Doe | `8f2a3c4d-1111-4a5b-9c6d-000000000001` | 1000 |
| Jane Doe | `8f2a3c4d-2222-4a5b-9c6d-000000000002` | 1000 |

### Kafka

Control Center: http://localhost:9021 → Topics → `balances` → Messages.
Consumer group do Balances: `balances`.

## Testes automatizados

```bash
go test ./...                  # ms-balance (requer gcc, por causa do CGO/sqlite3)
(cd wallet-core && go test ./...)
```

## Reiniciar do zero

```bash
docker compose down -v
docker compose up -d --build
```

Os bancos não usam volume persistente: cada recriação volta ao estado dos seeds.