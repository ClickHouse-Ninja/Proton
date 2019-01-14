# Basic reports

## By script name
```sql
SELECT
    dictGetString('Proton', 'Value', ScriptNameID) AS ScriptName
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        ScriptNameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY ScriptNameID
);
```

## By server name
```sql
SELECT
    dictGetString('Proton', 'Value', ServerNameID) AS ServerName
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        ServerNameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY ServerNameID
);
```

## By hostname
```sql
SELECT
    dictGetString('Proton', 'Value', HostnameID) AS Hostname
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        HostnameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY HostnameID
);
```

## By server and script
```sql
SELECT
    dictGetString('Proton', 'Value', ServerNameID) AS ServerName
    , dictGetString('Proton', 'Value', ScriptNameID) AS ScriptName
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        ServerNameID
        , ScriptNameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY ServerNameID, ScriptNameID
);
```

## By hostname and script
```sql
SELECT
    dictGetString('Proton', 'Value', HostnameID) AS Hostname
    , dictGetString('Proton', 'Value', ScriptNameID) AS ScriptName
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        HostnameID
        , ScriptNameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY HostnameID, ScriptNameID
);
```

## By hostname and server
```sql
SELECT
    dictGetString('Proton', 'Value', HostnameID) AS Hostname
    , dictGetString('Proton', 'Value', ServerNameID) AS ServerName
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        HostnameID
        , ServerNameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY HostnameID, ServerNameID
);
```

## By hostname, server and script
```sql
SELECT
    dictGetString('Proton', 'Value', HostnameID) AS Hostname
    , dictGetString('Proton', 'Value', ServerNameID) AS ServerName
    , dictGetString('Proton', 'Value', ScriptNameID) AS ScriptName
    , RequestCount
    , RequestTimeAvg
    , RequestTimeMax
    , RequestTimeQuantiles[1] AS RequestTime_90
    , RequestTimeQuantiles[2] AS RequestTime_95
    , RequestTimeQuantiles[3] AS RequestTime_99
    , UtimeTotal
    , StimeTotal
    , formatReadableSize(TrafficTotal) AS TrafficTotal
    , formatReadableSize(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM (
    SELECT
        HostnameID
        , ServerNameID
        , ScriptNameID
        , SUM(RequestCount)                    AS RequestCount
        , maxMerge(RequestTimeMax)             AS RequestTimeMax
        , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
        , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
        , SUM(UtimeTotal)           AS UtimeTotal
        , SUM(StimeTotal)           AS StimeTotal
        , SUM(DocumentSizeTotal)    AS TrafficTotal
        , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
    FROM proton.base_report
    GROUP BY HostnameID, ServerNameID, ScriptNameID
);
```