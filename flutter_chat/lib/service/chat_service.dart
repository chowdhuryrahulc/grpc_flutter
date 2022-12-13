// ignore_for_file: prefer_const_constructors

import 'dart:convert';

import 'package:crypto/crypto.dart';
import 'package:grpc/grpc.dart';

import '../gen/service.pbgrpc.dart';

class ChatService {
  // logic we want to use to interface with our server
  User user = User(); // the user that is logging in
  static BroadcastClient? client; // client used to communicate with the server
  // Static means if we set it up once, it will be set forever.



//! Constructor
  ChatService(String name) {
    // username enters this constructor
    //? THIS IS HOW WE PUT STUFF INSIDE GRPC MESSAGES
    user // user structure (object)
      ..clearName() // deletes name, in case a prior name is present
      ..name = name // set the incomming name to user name
      ..clearId() // detetes id
      ..id =
          sha256.convert(utf8.encode(user.name)).toString(); // generate new id
    // utf8 converts string into list of integers.
    // Then we use that list of integers and convert it to sha256
    // then convert sha256 hash to string

    client = BroadcastClient(ClientChannel(
        // ClientChannel has all information to connect to server
        "192.168.29.93", // ip address (in this case, server is running on localhost, client on phone, we use 10.0.2.2 by default)
        // (10.0.2.2 ip helps emulator access local host)
        
        //// 192.168.1.102, 192.168.56.1(IP4), 172.31.0.1(IP4) also works
        //// 0.0.0.0, 127.0.0.1, fe80::566c:1fbd:dda4:ae09%51(IP6) failed
        //// ngrok tried and failed? try again?
        //! ERROR SOLVED: 192.168.76.152 WORKED ❤ By trying all the IP4 ddresses after typing ipconfig in terminal
        //todo PROBLEMS REMAINING:
        //*1)2022/12/10 18:15:47 ERROR: Error with Stream: %!s(*gen.Broadcast_CreateStreamServer=0xc000282210) - Error: rpc error: code = Unavailable desc = transport is closing
        //*2)Socket connect host something. Might be solved using keepAlive
        //*3)Cretentials add (later)
        // Also controller and submit of MessagePage are working great
        //! ERROR2 Solved: ip changes everytime, thatswhy broadcst client doesnt work. SOL: just type ipconfig again, and put a new ip
        
        port: 8080,
        options: ChannelOptions(
          credentials: ChannelCredentials.insecure(), 
          codecRegistry: CodecRegistry(codecs: const [GzipCodec(), IdentityCodec()]),
          )));
          // userAgent: ??
          // backoffStrategy: ??
          
    // log(client.toString());    // This executes when we press the submit button
  }

  //! Methords to call to RPC Methords
  Future<Close> sendMessage(String body) async {
    // This will take in message and return a rpc Close type.
    //// var keepalive = KeepAlive(keepAlive: keepAlive, child: child)
    return client!.broadcastMessage(
      //todo This takes CallOptions. Means compression etc
      Message()
        ..id = user.id
        ..content = body
        ..timestamp = DateTime.now().toIso8601String(),
      options: CallOptions(compression: GzipCodec())
      // ? Can also be IdentityCodec().
      //todo Whats the difference btw both?
    );
  }

  // this will emit stream of message type (thats why async*)
  Stream<Message> recieveMessage() async* {
    Connect connect = Connect()
      ..user = user
      ..active = true;

    await for (var msg in client!.createStream(connect)) {
      // client!.createStream(connect) is the 2nd rpc function.
      // Takes connect, returns stream of messages
      yield msg;
    }
  }
}
