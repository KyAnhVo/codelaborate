#!/usr/bin/env python3
"""
Test script for Codelaborate server update message propagation.
Tests with 2 clients: one sends updates, the other receives them, then they swap roles.
"""

import socket
import struct
import time
import threading
import sys

HOST = 'localhost'
PORT = 8080
ROOM = int(sys.argv[1])

def recv_exact(sock, n):
    """Receive exactly n bytes from socket"""
    data = b''
    while len(data) < n:
        chunk = sock.recv(n - len(data))
        if not chunk:
            raise ConnectionError("Socket closed while reading")
        data += chunk
    return data

def send_create_room(sock):
    """Send room creation message"""
    # Operation 'C' (create), roomID = 0 (ignored for create)
    msg = struct.pack('!BI', ord('C'), 0)
    sock.sendall(msg)
    print("Sent: CREATE room")

def send_join_room(sock, room_id):
    """Send room join message"""
    # Operation 'J' (join), roomID
    msg = struct.pack('!BI', ord('J'), room_id)
    sock.sendall(msg)
    print(f"Sent: JOIN room {room_id}")

def send_update(sock, cursor_pos, delete_len, insert_str):
    """Send an update message"""
    # Message type: 0 = UPDATE
    # cursor_pos, delete_len, insert_len (uint64)
    # insert_str (string)
    insert_len = len(insert_str)
    
    msg = struct.pack('!BQQQ', 0, cursor_pos, delete_len, insert_len)
    msg += insert_str.encode('utf-8')
    
    sock.sendall(msg)
    print(f"Sent UPDATE: cursor={cursor_pos}, delete={delete_len}, insert='{insert_str}'")

def receive_update(sock):
    """Receive an update message"""
    try:
        # Read message type
        msg_type_bytes = recv_exact(sock, 1)
        msg_type = struct.unpack('!B', msg_type_bytes)[0]
        
        if msg_type == 1:  # CLOSECONN
            print("Received: CLOSECONN")
            return None
        
        # Read cursor_pos, delete_len, insert_len (3 x uint64 = 24 bytes)
        header = recv_exact(sock, 24)
        cursor_pos, delete_len, insert_len = struct.unpack('!QQQ', header)
        
        # Read insert string
        insert_str = ''
        if insert_len > 0:
            insert_bytes = recv_exact(sock, insert_len)
            insert_str = insert_bytes.decode('utf-8')
        
        print(f"Received UPDATE: cursor={cursor_pos}, delete={delete_len}, insert='{insert_str}'")
        
        return {
            'cursor_pos': cursor_pos,
            'delete_len': delete_len,
            'insert_str': insert_str
        }
    except ConnectionError as e:
        print(f"Connection error while receiving: {e}")
        return None
    except Exception as e:
        print(f"Error receiving update: {e}")
        return None

def client_receiver(sock, name, expected_count):
    """Receive updates from the server"""
    print(f"\n[{name}] Starting receiver, expecting {expected_count} messages...")
    received = []
    
    for i in range(expected_count):
        update = receive_update(sock)
        if update:
            received.append(update)
        else:
            break
    
    print(f"[{name}] Receiver done, got {len(received)} messages")
    return received

def run_test():
    """Run the full test sequence"""
    print("=" * 60)
    print("Starting Codelaborate Server Test")
    print("=" * 60)
    
    # Create two client sockets
    client1 = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client2 = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    
    try:
        # Connect both clients
        print("\n--- Connecting clients ---")
        client1.connect((HOST, PORT))
        print("Client 1 connected")
        time.sleep(0.1)
        
        client2.connect((HOST, PORT))
        print("Client 2 connected")
        time.sleep(0.1)
        
        # Client 1 creates room
        print("\n--- Client 1: Creating room ---")
        send_create_room(client1)
        time.sleep(0.2)
        
        # Client 2 joins room 0
        print("\n--- Client 2: Joining room 0 ---")
        send_join_room(client2, ROOM)
        time.sleep(0.2)
        
        # Test 1: Client 1 sends, Client 2 receives
        print("\n" + "=" * 60)
        print("TEST 1: Client 1 sends updates -> Client 2 receives")
        print("=" * 60)
        
        # Start client 2 receiver in background
        receiver_thread = threading.Thread(
            target=client_receiver, 
            args=(client2, "Client 2", 2)
        )
        receiver_thread.start()
        
        time.sleep(0.1)
        
        # Client 1 sends two updates
        send_update(client1, 0, 0, "Hello")
        time.sleep(0.2)
        send_update(client1, 5, 0, " World")
        time.sleep(0.2)
        
        receiver_thread.join(timeout=2)
        
        # Test 2: Client 2 sends, Client 1 receives
        print("\n" + "=" * 60)
        print("TEST 2: Client 2 sends updates -> Client 1 receives")
        print("=" * 60)
        
        # Start client 1 receiver in background
        receiver_thread = threading.Thread(
            target=client_receiver,
            args=(client1, "Client 1", 2)
        )
        receiver_thread.start()
        
        time.sleep(0.1)
        
        # Client 2 sends two updates
        send_update(client2, 11, 0, "!")
        time.sleep(0.2)
        send_update(client2, 6, 5, "Python")
        time.sleep(0.2)
        
        receiver_thread.join(timeout=2)
        
        # Test 3: Alternating updates
        print("\n" + "=" * 60)
        print("TEST 3: Alternating updates")
        print("=" * 60)
        
        # Start both receivers
        receiver1_thread = threading.Thread(
            target=client_receiver,
            args=(client1, "Client 1", 2)
        )
        receiver2_thread = threading.Thread(
            target=client_receiver,
            args=(client2, "Client 2", 2)
        )
        
        receiver1_thread.start()
        receiver2_thread.start()
        
        time.sleep(0.1)
        
        # Alternating sends
        send_update(client1, 0, 0, "A")
        time.sleep(0.2)
        send_update(client2, 1, 0, "B")
        time.sleep(0.2)
        send_update(client1, 2, 0, "C")
        time.sleep(0.2)
        send_update(client2, 3, 0, "D")
        time.sleep(0.2)
        
        receiver1_thread.join(timeout=2)
        receiver2_thread.join(timeout=2)
        
        print("\n" + "=" * 60)
        print("Test completed!")
        print("=" * 60)
        
    except ConnectionRefusedError:
        print("\nERROR: Could not connect to server. Is it running?")
        print(f"Make sure server is running on {HOST}:{PORT}")
    except Exception as e:
        print(f"\nERROR: {e}")
        import traceback
        traceback.print_exc()
    finally:
        print("\n--- Closing connections ---")
        client1.close()
        client2.close()

if __name__ == "__main__":
    run_test()
