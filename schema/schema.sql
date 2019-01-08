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
PARTITION BY toMonday(Timestamp)
ORDER     BY (
    Timestamp
);