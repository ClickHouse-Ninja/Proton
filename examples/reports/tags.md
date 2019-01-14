# Tags reports

## Example "Platform report"

## TABLE

```sql
CREATE TABLE IF NOT EXISTS proton.platform_report (
    Status                          Int16
    , OS                            String
    , Device                        String
    , Browser                       String
    , RequestCount                  UInt32
    , RequestTimeTotal              Float32
    , DocumentSizeTotal             UInt32
    , MemoryPeakTotal               UInt32
    , MemoryFootprintTotal          UInt32
    , UtimeTotal                    Float32
    , StimeTotal                    Float32

    , RequestTimeMax                AggregateFunction(Max, Float32)
    , DocumentSizeMax               AggregateFunction(Max, UInt32)
    , MemoryPeakMax                 AggregateFunction(Max, UInt32)
    , MemoryFootprintMax            AggregateFunction(Max, UInt32)
    , UtimeTotalMax                 AggregateFunction(Max, Float32)
    , StimeTotalMax                 AggregateFunction(Max, Float32)

    , RequestTimeTotalQuantiles     AggregateFunction(quantiles(0.9, 0.95, 0.99), Float32)
    , DocumentSizeTotalQuantiles    AggregateFunction(quantiles(0.9, 0.95, 0.99), UInt32)
    , MemoryPeakTotalQuantiles      AggregateFunction(quantiles(0.9, 0.95, 0.99), UInt32)
    , MemoryFootprintTotalQuantiles AggregateFunction(quantiles(0.9, 0.95, 0.99), UInt32)
    , UtimeTotalQuantiles           AggregateFunction(quantiles(0.9, 0.95, 0.99), Float32)
    , StimeTotalQuantiles           AggregateFunction(quantiles(0.9, 0.95, 0.99), Float32)
    , Timestamp                     DateTime
) Engine SummingMergeTree
PARTITION BY toYYYYMM(Timestamp)
ORDER BY (
     Status
    , OS
    , Device
    , Browser
    , Timestamp
);
```

## MATERIALIZED VIEW

```sql
CREATE MATERIALIZED VIEW proton.v_by_platform TO proton.platform_report AS
    SELECT
        Status
        , Tags.Value[indexOf(Tags.Name, 'OS')]      AS OS -- You can use cityHash64 like `cityHash64(Tags.Value[indexOf(Tags.Name, 'A')] AS OsID)`
        , Tags.Value[indexOf(Tags.Name, 'Device')]  AS Device
        , Tags.Value[indexOf(Tags.Name, 'Browser')] AS Browser
        , CAST(1 AS UInt32)      AS RequestCount
        , RequestTime            AS RequestTimeTotal
        , DocumentSize           AS DocumentSizeTotal
        , MemoryPeak             AS MemoryPeakTotal
        , MemoryFootprint        AS MemoryFootprintTotal
        , Utime                  AS UtimeTotal
        , Stime                  AS StimeTotal
        , arrayReduce('maxState', [RequestTime])     AS RequestTimeMax
        , arrayReduce('maxState', [DocumentSize])    AS DocumentSizeMax
        , arrayReduce('maxState', [MemoryPeak])      AS MemoryPeakMax
        , arrayReduce('maxState', [MemoryFootprint]) AS MemoryFootprintMax
        , arrayReduce('maxState', [Utime])           AS UtimeTotalMax
        , arrayReduce('maxState', [Stime])           AS StimeTotalMax
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [RequestTime])     AS RequestTimeTotalQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [DocumentSize])    AS DocumentSizeTotalQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [MemoryPeak])      AS MemoryPeakTotalQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [MemoryFootprint]) AS MemoryFootprintTotalQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [Utime])           AS UtimeTotalQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [Stime])           AS StimeTotalQuantiles
        , toStartOfMinute(Timestamp) AS Timestamp
        /* ^^^ YOU CAN CHANGE IT.
         * Example:
         *   toStartOfFiveMinute(Timestamp)                 - round up to 5 minutes
         *   toDateTime(intDiv(toUInt32(Timestamp), 5) * 5) - round up to 5 seconds
         */
    FROM proton.requests;
```

## REPORTS

## By OS
```sql
SELECT
    OS
    , SUM(RequestCount)                    AS RequestCount
    , maxMerge(RequestTimeMax)             AS RequestTimeMax
    , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
    , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
    , SUM(UtimeTotal)           AS UtimeTotal
    , SUM(StimeTotal)           AS StimeTotal
    , SUM(DocumentSizeTotal)    AS TrafficTotal
    , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM proton.platform_report
GROUP BY OS;
```

## By device
```sql
SELECT
    Device
    , SUM(RequestCount)                    AS RequestCount
    , maxMerge(RequestTimeMax)             AS RequestTimeMax
    , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
    , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
    , SUM(UtimeTotal)           AS UtimeTotal
    , SUM(StimeTotal)           AS StimeTotal
    , SUM(DocumentSizeTotal)    AS TrafficTotal
    , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM proton.platform_report
GROUP BY Device;
```


## By browser
```sql
SELECT
    Browser
    , SUM(RequestCount)                    AS RequestCount
    , maxMerge(RequestTimeMax)             AS RequestTimeMax
    , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
    , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
    , SUM(UtimeTotal)           AS UtimeTotal
    , SUM(StimeTotal)           AS StimeTotal
    , SUM(DocumentSizeTotal)    AS TrafficTotal
    , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM proton.platform_report
GROUP BY Browser;
```

## By device and browser
```sql
SELECT
    Device
    , Browser
    , SUM(RequestCount)                    AS RequestCount
    , maxMerge(RequestTimeMax)             AS RequestTimeMax
    , SUM(RequestTimeTotal) / RequestCount AS RequestTimeAvg
    , quantilesMerge(0.9,0.95,0.99)(RequestTimeTotalQuantiles) AS RequestTimeQuantiles
    , SUM(UtimeTotal)           AS UtimeTotal
    , SUM(StimeTotal)           AS StimeTotal
    , SUM(DocumentSizeTotal)    AS TrafficTotal
    , SUM(MemoryFootprintTotal) AS MemoryFootprintTotal
FROM proton.platform_report
GROUP BY Device, Browser;
```