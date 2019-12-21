import hazelcast
import logging
from hazelcast.serialization.api import Portable


class Reading(Portable):
    def __init__(self, name, ts):
        self.Name = name
        self.Timestamp = ts


if __name__ == "__main__":
    logging.basicConfig()
    logging.getLogger().setLevel(logging.DEBUG)
    config = hazelcast.ClientConfig()
    config.network_config.addresses.append('192.168.99.100:5701')
    config.serialization_config.data_serializable_factories[1] = {1: Reading}

    client = hazelcast.HazelcastClient(config)
    myList = client.get_list("hazelcaster")
    print("List size:", myList.size().result())

    readings = []
    for r in myList.get_all().result():
        value = r.loads()
        readings.append(Reading(value["Name"], value["Timestamp"]))

    for r in readings:
        print(r.Timestamp, r.Name)
