CREATE DATABASE IF NOT EXISTS proton;

CREATE TABLE IF NOT EXISTS proton.requests (
    Hostname          String
    , Schema          String
    , Status          Int16
    , ServerName      String
    , ScriptName      String
    , RequestCount    UInt32
    , RequestTime     Float32
    , DocumentSize    UInt32
    , MemoryPeak      UInt32
    , MemoryFootprint UInt32
    , Utime           Float32
    , Stime           Float32
    , Tags Nested (
        Name    String
        , Value String
    )
    , Timers Nested (
        HitCount    UInt32
        , Value     Float32
        , Utime     Float32
        , Stime     Float32
        , TagsName  Array(String)
        , TagsValue Array(String)
    )
    , Timestamp DateTime
) Engine = MergeTree
PARTITION BY toYYYYMM(Timestamp)
ORDER     BY (
    Schema
    , toStartOfMinute(Timestamp)
);

CREATE TABLE IF NOT EXISTS proton.dictionary (
    ID       UInt64
    , Value  String
    , Column String
) Engine ReplacingMergeTree
PARTITION BY Column
ORDER     BY ID;