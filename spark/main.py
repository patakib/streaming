from pyspark.sql import SparkSession
from pyspark.sql.functions import *
from pyspark.sql.types import *
from pyspark.sql.functions import from_json

import os


os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.spark:spark-streaming-kafka-0-10_2.12:3.2.0,org.apache.spark:spark-sql-kafka-0-10_2.12:3.2.0 pyspark-shell'

spark = SparkSession \
    .builder \
    .appName("streamingWithSparkAndKafka") \
    .config("spark.streaming.stopGracefullyOnShutdown", True) \
    .master("local[*]") \
    .getOrCreate()
    
spark.sparkContext.setLogLevel("ERROR")

KAFKA_TOPIC_01 = "streamingserver.streaming.activity"
KAFKA_SERVER = "localhost:9092"

activity_schema = StructType([
        StructField("payload", StructType([
            StructField("after", StructType([
                StructField("id", IntegerType()),
                StructField("user_id", IntegerType()),
                StructField("activity_type", StringType()),
                StructField("intensity", IntegerType()),
                StructField("duration", IntegerType())
            ]), True)
        ]))
    ])

activities = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", KAFKA_SERVER) \
    .option("subscribe", KAFKA_TOPIC_01) \
    .option("startingOffsets", "earliest") \
    .load() \
    .withColumn("value", from_json(col("value").cast("string"), activity_schema)) \
    .select("value.payload.after.*")

activities.printSchema()

df = activities.writeStream.format("console").outputMode("append").start().awaitTermination()