package schedule

const reservationSQL = `SELECT seat.id AS seat_id FROM seat
INNER JOIN area ON area.id=seat."areaId" AND area."deletedAt" IS NULL
LEFT JOIN reservation ON reservation."seatId"=seat.id AND reservation."deletedAt" IS NULL
WHERE area."eventId"=$1 AND reservation."eventDateId"=$2`

const areaSQL = `SELECT id,label,svg,price FROM area
WHERE "eventId"=$1 AND "deletedAt" IS NULL`

const seatSQL = `SELECT seat.id,seat.cx,seat.cy,seat.row,seat.number,
seat."areaId","order"."userId" FROM seat
LEFT JOIN reservation ON reservation."seatId"=seat.id
AND reservation."canceledAt" IS NULL AND reservation."deletedAt" IS NULL
LEFT JOIN event_date AS eventDate ON reservation."eventDateId"=eventDate.id
LEFT JOIN "order" ON reservation."orderId"="order".id
AND "order"."canceledAt" IS NULL AND "order"."deletedAt" IS NULL
WHERE seat."areaId"=$1 AND (eventDate.id=$2 OR eventDate.id IS NULL)`
