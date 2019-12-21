import hazelcast
import logging
from hazelcast.serialization.api import Portable


class Reading(Portable):
    CLASS_ID = 1
    FACTORY_ID = 1

    def __init__(self, name=None, ts=None):
        self.Name = name
        self.Timestamp = ts

    def get_class_id(self):
        return self.CLASS_ID

    def get_factory_id(self):
        return self.FACTORY_ID

    def write_portable(self, writer):
        writer.write_utf("Name", self.Name)
        writer.write_long("Timestamp", self.Timestamp)

    def read_portable(self, reader):
        self.Name = reader.read_utf("Name")
        self.Timestamp = reader.read_long("Timestamp")


if __name__ == "__main__":
    logging.basicConfig()
    logging.getLogger().setLevel(logging.INFO)
    config = hazelcast.ClientConfig()
    config.network_config.addresses.append('192.168.99.100:5701')
    config.serialization_config.add_portable_factory(1, {1: Reading})

    client = hazelcast.HazelcastClient(config)
    myList = client.get_list("hazelcaster-test")
    print("List size:", myList.size().result())

    for item in myList.get_all().result():
        print(item.Timestamp, item.Name)

