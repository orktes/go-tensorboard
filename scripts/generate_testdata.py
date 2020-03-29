import tensorflow as tf

writer = tf.summary.create_file_writer("testdata/tfevents")

with writer.as_default():
  for step in range(100):
    # other model code would go here
    tf.summary.scalar("my_metric", 0.5*step, step=step)
    writer.flush()