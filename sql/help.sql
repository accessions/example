-- 多维度查询json
SELECT
    h.date,
    CONCAT('[',
           GROUP_CONCAT(concat('{"', a.name,'":', a.activate,'}'))
        ,']') as activate,
    (SELECT GROUP_CONCAT(CONCAT('"',name,'"')) FROM aiot_customer) as name
FROM
    (SELECT date_format(date_add('2021-01-14', INTERVAL - t.help_topic_id DAY), '%Y-%m-%d') AS 'date' FROM mysql.help_topic AS t WHERE t.help_topic_id <=  TO_DAYS('2021-01-14')-TO_DAYS('2021-01-05') ) AS h
        LEFT JOIN
    (SELECT a.`name`, b.date, SUM(activated_day) as activate FROM aiot_customer as a LEFT JOIN aiot_device_day_statistics as b ON a.id = b.org_id GROUP BY a.id, b.date) as a ON a.date = h.date
GROUP BY h.date

-- 最近30天
SELECT
    h.date,
    p.showCount,
    p.clickCount,
    p.count
FROM
    (SELECT date_format(date_add('2020-12-06', INTERVAL + t.help_topic_id DAY ), '%Y%m%d') AS 'date' FROM	mysql.help_topic AS t WHERE t.help_topic_id <= TO_DAYS('2021-01-05') - TO_DAYS('2020-12-06') ) AS h
        LEFT JOIN
    (SELECT
         DATE_FORMAT(a.reportTime,'%Y%m%d') tmp,
         SUM(a.showCount) showCount,
         SUM(a.clickCount) clickCount,
         SUM(c.terminalcount) count
     FROM  ad_plan_report as a LEFT JOIN ad_plan as b ON a.planId = b.id
         LEFT JOIN ad_plan_details as c ON c.planId=b.id
     WHERE a.reportTime BETWEEN '2020-12-06' AND '2021-01-05'  and b.id IS NOT NULL
     GROUP BY a.reportTime
    ) AS p ON p.tmp = h.date  ORDER BY h.date

-- 分组后比较

(SELECT ROUND(COUNT(1)/(SELECT COUNT(1) deviceCount
FROM aiot_device_activate a
LEFT JOIN aiot_device b ON a.serialNum=b.serialNum
LEFT JOIN airiskwhit c on b.orgId=c.UnitId AND c.Type=1
WHERE activatedTime BETWEEN date_format(DATE_SUB(curdate(), INTERVAL 6 MONTH),'%Y%m') AND CURDATE() AND c.UnitId IS NULL) * 100,2) onlineRate FROM (
SELECT a.serialNum
FROM aiot_connect_day_log a
LEFT JOIN airiskwhit b on a.customerId=b.UnitId AND b.Type=1
WHERE a.statDateNum BETWEEN DATE_FORMAT(CURDATE(),'%Y%m') AND DATE_FORMAT(CURDATE(),'%Y%m%d') AND b.UnitId IS NULL
GROUP BY a.serialNum HAVING count(1)>10
) s) as online_rate,

-- 同一张表多次联查
SELECT  a.date, a.a, b.b, c.c
FROM
    (SELECT  DATE_FORMAT(StatDate,'%H') as date, SUM(OnlineCount) as a FROM aiot_connect_hour_log WHERE DATE_FORMAT(StatDate,'%Y-%m-%d') = '{date}' GROUP BY date) as a
        LEFT JOIN
    (SELECT  DATE_FORMAT(StatDate,'%H') as date, SUM(OnlineCount) as b FROM aiot_connect_hour_log WHERE DATE_FORMAT(StatDate,'%Y-%m-%d') = SUBDATE('{date}', INTERVAL +1 day) GROUP BY date) as b ON a.date=b.date
        LEFT JOIN
    (SELECT  DATE_FORMAT(StatDate,'%H') as date, SUM(OnlineCount) as c FROM aiot_connect_hour_log WHERE DATE_FORMAT(StatDate,'%Y-%m-%d') = SUBDATE('{date}', INTERVAL +7 day) GROUP BY date) as c ON a.date=c.date";
