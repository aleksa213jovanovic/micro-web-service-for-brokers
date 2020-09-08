select count(*) from currencySQL.rates GROUP BY creation_date;

WITH normalized_first as
(select rate from currencySQL.rates 
where id1=13 and id2=2  and creation_date BETWEEN '2020-07-01' AND '2020-07-02'
union 
select 1/rate from currencySQL.rates 
where id2=13 and id1=2 and creation_date BETWEEN '2020-07-01' AND '2020-07-02')
select avg(rate) from normalized_first;

WITH normalized_first as 
(select rate from currencySQL.rates 
where id1=48 and id2=131  and creation_date BETWEEN '2020.07.23' AND '2020.07.26' 
union
select 1/rate from currencySQL.rates where id2=48 and id1=131 and creation_date BETWEEN '2020.07.23' AND '2020.07.26') 
select avg(rate) from normalized_first;

select * from currencySQL.rates;


select count(*) as n from currencySQL.rates r join currencySQL.currency c
on c.id = r.id1 or c.id = r.id2
group by c.id
having n!=1169;

select * from currencySQL.rates where  (id1=103 and id2=47 or id2=103 and id1=47);

select id from currencySQL.currency where name='EUR';
insert into currencySQL.rates (nick_name, email) VALUES('aleksa','ee@gmail.com');

select count(*) from currencySQL.rates;

select * from currencySQL.rates;

insert into currencySQL.rates (id1,id2,rate,creation_date) VALUES (1,2,1.4,'2020-08-12');


select id from currencySQL.currency where name='RSD';
select name from currencySQL.currency where id=129 or id=24;

select * from currencySQL.clients;

select avg(1.0000001)

select * from currencySQL.rates where id1=10 and id2=66 and creation_date='2020-07-27 00:00:00';

WITH normalized_first as 
(select * from currencySQL.rates where id1=72 and id2=19  and creation_date BETWEEN '2020-07-27 00:00:00' AND '2020-07-27 00:00:00'  
 union select * from currencySQL.rates where id1=66 and id2=10 and creation_date BETWEEN '2020-07-27 00:00:00'  AND '2020-07-27 00:00:00' )
 select avg(rate) from normalized_first;
 
 
 
 WITH normalized_first as 
	(select rate from currencySQL.rates where id1=89 and id2=117  and creation_date BETWEEN '2020-07-28 00:00:00' AND '2020-07-29 00:00:00'  
	union 
	select 1/rate from currencySQL.rates where id1=117 and id2=89 and creation_date BETWEEN '2020-07-28 00:00:00' AND '2020-07-29 00:00:00')
select avg(rate) from normalized_first;
