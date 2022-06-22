# ДЗ 03
## п. 1 Добавить параметры для фильтрации товаров по диапазону цены
Для path `/items` с помощью get запроса и параметров `price_min` и `price_max` можно передать параметры для фильтрации.
```swagger
  /items:
    get:
      tags:
      - Item
      summary: Lists Items with filters
      operationId: ListItems
      parameters:
      - name: price_min
        in: query
        description: Lower price limit
        required: false
        schema:
          type: integer
          format: int64
      - name: price_max
        in: query
        description: Upper price limit
        required: false
        schema:
          type: integer
          format: int64
      responses:
        200:
          description: successful operation
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Item'
        400:
          description: Invalid price range
          content: {}
```

## п. 2 Добавить в спецификацию объект Order (заказ), подумать, какие поля у нее должны быть и какие эндпоинты потребуется описать.
Добавлена схема `Order`, описывающая объект "Заказ" со своими полями.  
Добавлены эндпоинты `/orders` и `/orders/{orderId}`, позволяющие работать с заказом.

## п. 3 *Спроектировать REST API для упрощенной версии Twitter
В каталоге `lesson03/simpleTwitterAPI/api/swagger.yaml` расположен yaml-файл с описанием по стандарту OpenAPI v3.  
Также с помощью сервиса `https://editor.swagger.io/` был сгенерирован код серверной части 
(чисто для проверки функционала, т.к. не стояло задачи реализовать backend, то все функции - просто заглушки).

## Прогресс по курсовому проекту
Добавлена спецификация OpenAPI и сгенерирована заглушка для backend-сервера.  
https://github.com/ptsypyshev/shortlink/pull/3