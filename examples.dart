import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_pdfview/flutter_pdfview.dart';

void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      home: PDFScreen(),
    );
  }
}

class PDFScreen extends StatefulWidget {
  @override
  _PDFScreenState createState() => _PDFScreenState();
}

class _PDFScreenState extends State<PDFScreen> {
  late Future<Uint8List> pdfData;

  @override
  void initState() {
    super.initState();
    pdfData = fetchPDFData();
  }

  Future<Uint8List> fetchPDFData() async {
    final response = await http.get(Uri.parse('http://your_backend_ip:8080/getpdf'));
    if (response.statusCode == 200) {
      return response.bodyBytes;
    } else {
      throw Exception('Failed to fetch PDF');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('PDF Viewer'),
      ),
      body: FutureBuilder<Uint8List>(
        future: pdfData,
        builder: (context, snapshot) {
          if (snapshot.hasData) {
            return PDFView(
              filePath: null, // No need for a file path
              data: snapshot.data!,
            );
          } else if (snapshot.hasError) {
            return Center(child: Text('Error loading PDF'));
          }
          return Center(child: CircularProgressIndicator());
        },
      ),
    );
  }
}