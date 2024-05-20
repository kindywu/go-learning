-- 表结构
CREATE TABLE IF NOT EXISTS `UserRole` (
  `Id` bigint NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) COLLATE utf8mb4_bin NOT NULL,
  `Description` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL,
  `IsEnabled` bit(1) NOT NULL,
  `Created` date NOT NULL,
  `CreatedBy` varchar(200) COLLATE utf8mb4_bin NOT NULL,
  `Updated` date DEFAULT NULL,
  `UpdatedBy` varchar(200) COLLATE utf8mb4_bin DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- 准备数据
INSERT INTO `UserRole` (`Id`, `Name`, `Description`, `IsEnabled`, `Created`, `CreatedBy`, `Updated`, `UpdatedBy`) VALUES
	(1, 'R1', NULL, b'1', '2024-03-29', 'Admin', NULL, NULL),
	(2, 'R2', NULL, b'1', '2024-03-29', 'Admin', '2024-03-29', 'John Smith'),
	(3, 'R3', NULL, b'0', '2024-03-29', 'John SMITH', '2024-03-29', 'Ben SMITH'),
	(4, 'R4', NULL, b'1', '2024-03-29', 'bEn SMITH', '2024-03-29', 'BEN SMITH');

-- 写出SQL语句，查询结果是以下内容
UserName        NoOfCreatedRole  NoOfCreatedAndEnabledRole  NoOfUpdatedRole
JOHN SMITH      1                          -1                  1
BEN SMITH       1                           1                  2
ADMIN           2                           2                 -1


SELECT COALESCE(NULL, NULL, 1,2,3)

SELECT
 UserName,
 NoOfCreatedRoles,
 NoOfCreatedAndEnabledRoles, IFNULL(NoOfUpdatedRoles,-1) AS NoOfUpdatedRoles
FROM
 (
SELECT TRIM(UPPER(CreatedBy)) AS UserName
FROM
 UserRole UNION
SELECT TRIM(UPPER(UpdatedBy)) AS UserName
FROM
 UserRole
) AS T1
LEFT JOIN
 (
SELECT TRIM(UPPER(CreatedBy)) AS Uname, COUNT(CreatedBy) AS NoOfCreatedRoles, IFNULL(SUM(CASE WHEN
 IsEnabled = 1 THEN
 1 END
), - 1) AS NoOfCreatedAndEnabledRoles
FROM
 UserRole
GROUP BY
 Uname
) AS T2 ON T1.UserName = T2.Uname
LEFT JOIN
 (
SELECT TRIM(UPPER(UpdatedBy)) AS Uname, IFNULL(COUNT(UpdatedBy), - 1) AS NoOfUpdatedRoles
FROM
 UserRole
GROUP BY
 Uname
) AS T3 ON T1.UserName = T3.Uname
WHERE
 UserName IS NOT NULL
ORDER BY
 UserName DESC;


 SELECT t1.UserName,t1.NoOfCreatedRole,t1.NoOfCreatedAndEnabledRole,COALESCE(t2.NoOfUpdatedRole,-1) AS NoOfUpdatedRole
FROM 
(SELECT 
 UPPER(CreatedBy) AS UserName,
 COUNT(ID) AS NoOfCreatedRole,
 SUM(CASE WHEN IsEnabled=1 THEN 1 ELSE -1 END) AS NoOfCreatedAndEnabledRole
FROM UserRole 
GROUP BY UPPER(CreatedBy)) AS t1
LEFT JOIN 
(SELECT 
 UPPER(UpdatedBy) AS UserName,
 COUNT(CASE WHEN IsEnabled=1 THEN 1 ELSE -1 END) AS NoOfUpdatedRole
FROM UserRole 
WHERE UpdatedBy IS NOT NULL 
GROUP BY UPPER(UpdatedBy)) AS t2
ON t1.UserName = t2.UserName
ORDER BY UserName DESC 
