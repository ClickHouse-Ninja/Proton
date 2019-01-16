CREATE TABLE IF NOT EXISTS proton.base_report (
    Schema                     String
    , Status                   Int16
    , HostnameID               UInt64
    , ServerNameID             UInt64
    , ScriptNameID             UInt64
    , RequestCount             UInt32
    , RequestTimeTotal         Float32
    , DocumentSizeTotal        UInt32
    , MemoryPeakTotal          UInt32
    , MemoryFootprintTotal     UInt32
    , UtimeTotal               Float32
    , StimeTotal               Float32

    , RequestTimeMax           AggregateFunction(Max, Float32)
    , DocumentSizeMax          AggregateFunction(Max, UInt32)
    , MemoryPeakMax            AggregateFunction(Max, UInt32)
    , MemoryFootprintMax       AggregateFunction(Max, UInt32)
    , UtimeMax                 AggregateFunction(Max, Float32)
    , StimeMax                 AggregateFunction(Max, Float32)

    , RequestTimeQuantiles     AggregateFunction(quantiles(0.9, 0.95, 0.99), Float32)
    , DocumentSizeQuantiles    AggregateFunction(quantiles(0.9, 0.95, 0.99), UInt32)
    , MemoryPeakQuantiles      AggregateFunction(quantiles(0.9, 0.95, 0.99), UInt32)
    , MemoryFootprintQuantiles AggregateFunction(quantiles(0.9, 0.95, 0.99), UInt32)
    , UtimeQuantiles           AggregateFunction(quantiles(0.9, 0.95, 0.99), Float32)
    , StimeQuantiles           AggregateFunction(quantiles(0.9, 0.95, 0.99), Float32)
    , Timestamp                DateTime
) Engine SummingMergeTree
PARTITION BY toYYYYMM(Timestamp)
ORDER BY (
    Schema
    , Timestamp
    , Status
    , HostnameID
    , ServerNameID
    , ScriptNameID
);

CREATE MATERIALIZED VIEW proton.v_base_report TO proton.base_report AS
    SELECT
        Schema
        , Status
        , cityHash64(Hostname)   AS HostnameID
        , cityHash64(ServerName) AS ServerNameID
        , cityHash64(ScriptName) AS ScriptNameID
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
        , arrayReduce('maxState', [Utime])           AS UtimeMax
        , arrayReduce('maxState', [Stime])           AS StimeMax
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [RequestTime])     AS RequestTimeQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [DocumentSize])    AS DocumentSizeQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [MemoryPeak])      AS MemoryPeakQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [MemoryFootprint]) AS MemoryFootprintQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [Utime])           AS UtimeQuantiles
        , arrayReduce('quantilesState(0.90,0.95,0.99)', [Stime])           AS StimeQuantiles
        , toDateTime(intDiv(toUInt32(Timestamp), 15) * 15) AS Timestamp
        /* ^^^ YOU CAN CHANGE IT.
         * Example:
         *   toStartOfMinute(Timestamp)                     - round up to 1 minute
         *   toStartOfFiveMinute(Timestamp)                 - round up to 5 minutes
         *   toDateTime(intDiv(toUInt32(Timestamp), 5) * 5) - round up to 5 seconds
         */
    FROM proton.requests;