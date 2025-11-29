#!/usr/bin/env python3
"""
Debug script to test server protocol compliance.
Mimics Qt client behavior exactly.
"""

import socket
import struct
import time

HOST = 'localhost'
PORT = 8080

def recv_exact(sock, n):
    """Receive exactly n bytes from socket"""
    data = b''
    while len(data) < n:
        chunk = sock.recv(n - len(data))
        if not chunk:
            raise ConnectionError(f"Socket closed while reading (got {len(data)}/{n} bytes)")
        data += chunk
    return data

def test_create_room():
    """Test room creation flow"""
    print("=" * 60)
    print("TEST: Create Room")
    print("=" * 60)
    
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((HOST, PORT))
    print("✓ Connected to server")
    
    # Send CREATE message
    msg = struct.pack('!BI', ord('C'), 0)
    sock.sendall(msg)
    print(f"✓ Sent CREATE: {msg.hex()}")
    
    # Receive entry response
    status_byte = recv_exact(sock, 1)
    status = struct.unpack('!B', status_byte)[0]
    print(f"✓ Received status byte: {status} (0x{status:02x})")
    
    if status == 2:  # ENTRY_OK
        print("✓ Status = ENTRY_OK (2)")
        room_id_bytes = recv_exact(sock, 4)
        room_id = struct.unpack('!I', room_id_bytes)[0]
        print(f"✓ Received room ID: {room_id}")
        print(f"✓ Room ID bytes: {room_id_bytes.hex()}")
    elif status == 3:  # ENTRY_ERR
        print("✗ Status = ENTRY_ERR (3)")
        sock.close()
        return None
    else:
        print(f"✗ Unknown status: {status}")
        sock.close()
        return None
    
    print(f"\n✓ Room created successfully: Room {room_id}")
    return sock, room_id

def test_join_room(room_id):
    """Test joining existing room"""
    print("\n" + "=" * 60)
    print(f"TEST: Join Room {room_id}")
    print("=" * 60)
    
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((HOST, PORT))
    print("✓ Connected to server")
    
    # Send JOIN message
    msg = struct.pack('!BI', ord('J'), room_id)
    sock.sendall(msg)
    print(f"✓ Sent JOIN room {room_id}: {msg.hex()}")
    
    # Receive entry response
    status_byte = recv_exact(sock, 1)
    status = struct.unpack('!B', status_byte)[0]
    print(f"✓ Received status byte: {status} (0x{status:02x})")
    
    if status == 2:  # ENTRY_OK
        print("✓ Status = ENTRY_OK (2)")
        returned_room_id_bytes = recv_exact(sock, 4)
        returned_room_id = struct.unpack('!I', returned_room_id_bytes)[0]
        print(f"✓ Received room ID: {returned_room_id}")
        
        if returned_room_id != room_id:
            print(f"⚠ Warning: Requested room {room_id}, got {returned_room_id}")
    else:
        print(f"✗ Failed to join room")
        sock.close()
        return None
    
    print(f"\n✓ Joined room successfully")
    return sock

def test_send_update(sock, cursor_pos, delete_len, insert_str):
    """Test sending an update message"""
    print(f"\n  → Sending update: pos={cursor_pos}, del={delete_len}, ins='{insert_str}'")
    
    insert_bytes = insert_str.encode('utf-8')
    insert_len = len(insert_bytes)
    
    # Message: [op:1][cursor:8][delete:8][insert:8][string:N]
    msg = struct.pack('!BQQQ', 0, cursor_pos, delete_len, insert_len)
    msg += insert_bytes
    
    sock.sendall(msg)
    print(f"    Sent {len(msg)} bytes: op=0(UPDATE), cursor={cursor_pos}, del={delete_len}, ins={insert_len}")
    print(f"    Hex: {msg.hex()}")

def test_recv_update(sock, timeout=2):
    """Test receiving an update message"""
    sock.settimeout(timeout)
    
    try:
        # Read message type
        msg_type_byte = recv_exact(sock, 1)
        msg_type = struct.unpack('!B', msg_type_byte)[0]
        print(f"\n  ← Received message type: {msg_type} (0x{msg_type:02x})")
        
        if msg_type == 1:  # CLOSECONN
            print("    Type: CLOSECONN")
            # Still need to read the rest
            cursor_pos = struct.unpack('!Q', recv_exact(sock, 8))[0]
            delete_len = struct.unpack('!Q', recv_exact(sock, 8))[0]
            insert_len = struct.unpack('!Q', recv_exact(sock, 8))[0]
            insert_str = recv_exact(sock, insert_len).decode('utf-8') if insert_len > 0 else ""
            print(f"    (CLOSECONN data: pos={cursor_pos}, del={delete_len}, ins={insert_len})")
            return None
        elif msg_type == 0:  # UPDATE
            print("    Type: UPDATE")
        else:
            print(f"    ✗ Unknown message type: {msg_type}")
            return None
        
        # Read update data
        cursor_pos = struct.unpack('!Q', recv_exact(sock, 8))[0]
        delete_len = struct.unpack('!Q', recv_exact(sock, 8))[0]
        insert_len = struct.unpack('!Q', recv_exact(sock, 8))[0]
        
        insert_str = ""
        if insert_len > 0:
            insert_bytes = recv_exact(sock, insert_len)
            insert_str = insert_bytes.decode('utf-8')
        
        print(f"    pos={cursor_pos}, del={delete_len}, ins={insert_len}, text='{insert_str}'")
        print(f"    Text hex: {insert_str.encode('utf-8').hex()}")
        
        return {
            'cursor_pos': cursor_pos,
            'delete_len': delete_len,
            'insert_len': insert_len,
            'insert_str': insert_str
        }
    except socket.timeout:
        print(f"  ✗ Timeout waiting for update (waited {timeout}s)")
        return None
    except Exception as e:
        print(f"  ✗ Error receiving update: {e}")
        import traceback
        traceback.print_exc()
        return None

def test_two_clients_collaboration():
    """Test two clients collaborating"""
    print("\n" + "=" * 60)
    print("TEST: Two Client Collaboration")
    print("=" * 60)
    
    # Client 1 creates room
    sock1, room_id = test_create_room()
    if not sock1:
        return
    
    time.sleep(0.2)
    
    # Client 2 joins room
    sock2 = test_join_room(room_id)
    if not sock2:
        sock1.close()
        return
    
    time.sleep(0.2)
    
    # Client 1 sends "Hello"
    print("\n--- Client 1 sends 'Hello' ---")
    test_send_update(sock1, 0, 0, "Hello")
    time.sleep(0.2)
    
    # Client 2 should receive it
    print("\n--- Client 2 receives ---")
    update = test_recv_update(sock2)
    if update and update['insert_str'] == "Hello":
        print("✓ Client 2 received correct update!")
    else:
        print("✗ Client 2 did not receive correct update")
    
    time.sleep(0.2)
    
    # Client 2 sends " World"
    print("\n--- Client 2 sends ' World' ---")
    test_send_update(sock2, 5, 0, " World")
    time.sleep(0.2)
    
    # Client 1 should receive it
    print("\n--- Client 1 receives ---")
    update = test_recv_update(sock1)
    if update and update['insert_str'] == " World":
        print("✓ Client 1 received correct update!")
    else:
        print("✗ Client 1 did not receive correct update")
    
    sock1.close()
    sock2.close()
    print("\n✓ Collaboration test complete")

def test_check_for_extra_data_after_join():
    """Check if server sends unexpected data after join"""
    print("\n" + "=" * 60)
    print("TEST: Check for Extra Data After Join")
    print("=" * 60)
    
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((HOST, PORT))
    
    # Create room
    msg = struct.pack('!BI', ord('C'), 0)
    sock.sendall(msg)
    
    # Read entry response
    status = struct.unpack('!B', recv_exact(sock, 1))[0]
    room_id = struct.unpack('!I', recv_exact(sock, 4))[0]
    print(f"✓ Joined room {room_id}")
    
    # Check if there's any extra data
    sock.settimeout(0.5)
    try:
        extra = sock.recv(1)
        if extra:
            print(f"✗ Server sent extra data after join: {extra.hex()}")
            print(f"   This byte ({int.from_bytes(extra, 'big')}) will be interpreted as message type!")
        else:
            print("✓ No extra data after join")
    except socket.timeout:
        print("✓ No extra data after join (timeout)")
    
    sock.close()

def main():
    print("Codelaborate Server Debug Script")
    print("=" * 60)
    
    try:
        # Test 1: Basic entry
        test_check_for_extra_data_after_join()
        
        # Test 2: Two clients
        test_two_clients_collaboration()
        
        print("\n" + "=" * 60)
        print("All tests complete!")
        print("=" * 60)
        
    except ConnectionRefusedError:
        print("\n✗ ERROR: Could not connect to server")
        print(f"  Make sure server is running on {HOST}:{PORT}")
    except Exception as e:
        print(f"\n✗ ERROR: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    main()
