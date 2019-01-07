CREATE DATABASE IF NOT EXISTS proton;

CREATE TABLE IF NOT EXISTS proton.request (
    Hostname          String
    , Schema          String
    , Status          Int16
    , ServerName      String
    , ScriptName      String
    , RequestCount    Int64
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
PARTITION BY tuple()
ORDER     BY tuple();