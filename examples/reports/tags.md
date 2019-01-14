# Tags reports

Example `MATERIALIZED VIEW`

```sql
CREATE MATERIALIZED VIEW proton.v_by_os_and_browser TO proton.by_os_and_browser_report AS
    SELECT
        Schema
        , Status
        , Tags.Value[indexOf(Tags.Name, 'OS')]      AS OS -- You can use cityHash64 like `cityHash64(Tags.Value[indexOf(Tags.Name, 'A')] AS OsID)`
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