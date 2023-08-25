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

KAFKA_TOPIC_01 = "streamingdbserver.streaming.activity"
KAFKA_SERVER = "localhost:9092"

activities = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", KAFKA_SERVER) \
    .option("subscribe", KAFKA_TOPIC_01) \
    .option("startingOffsets", "earliest") \
    .load()

activity_schema = StructType([
            StructField("id", IntegerType()),
            StructField("user_id", IntegerType()),
            StructField("activity_type", StringType()),
            StructField("intensity", IntegerType()),
            StructField("duration", IntegerType())
        ])

json_df = activities.selectExpr("cast(value as string) as value")
json_expanded_df = json_df.withColumn("value", from_json(json_df["value"], activity_schema)).select("value.*")
json_expanded_df.printSchema()

selected_df = json_expanded_df.select("id", "user_id", "activity_type", "intensity", "duration")

df = selected_df.writeStream.format("console").start().awaitTermination()