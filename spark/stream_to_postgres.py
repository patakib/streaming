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
    .config("spark.jars", "./postgresql-42.6.0.jar") \
    .master("local[*]") \
    .getOrCreate()
    
spark.sparkContext.setLogLevel("ERROR")

KAFKA_TOPIC_01 = "streamingserver.streaming.user"
KAFKA_TOPIC_02 = "streamingserver.streaming.activity"
KAFKA_SERVER = "localhost:9092"

user_schema = StructType([
        StructField("payload", StructType([
            StructField("after", StructType([
                StructField("id", IntegerType()),
                StructField("birth_year", IntegerType()),
                StructField("location", StringType()),
                StructField("gender", StringType())
            ]), True)
        ]))
    ])

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

users = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", KAFKA_SERVER) \
    .option("subscribe", KAFKA_TOPIC_01) \
    .option("startingOffsets", "earliest") \
    .load() \
    .withColumn("value", from_json(col("value").cast("string"), user_schema)) \
    .select("value.payload.after.*")

activities = spark.readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", KAFKA_SERVER) \
    .option("subscribe", KAFKA_TOPIC_02) \
    .option("startingOffsets", "earliest") \
    .load() \
    .withColumn("value", from_json(col("value").cast("string"), activity_schema)) \
    .select("value.payload.after.*")
activities = activities.withColumnRenamed("id", "activity_id")

activities.printSchema()
users.printSchema()

df = users.join(activities, users.id == activities.user_id, 'inner')
df = df.select("user_id", "activity_id", "birth_year", "location", "gender", "activity_type", "intensity", "duration")

def write_stream_to_postgres(dataframe, epoch_id) -> None:
    dataframe.write \
    .mode("append") \
    .format("jdbc") \
    .option("url", "jdbc:postgresql://localhost:5431/streaming") \
    .option("driver", "org.postgresql.Driver") \
    .option("dbtable", "user_activity") \
    .option("user", "streamuser") \
    .option("password", "streampass") \
    .save()

df.writeStream \
    .foreachBatch(write_stream_to_postgres) \
    .start() \
    .awaitTermination()