<h1 align="center">pre-requisites: make sure you have go installed:</h2>


<h2 align="center">
    <code> go -v </code>
</h2>
<h2 align="center">or:</h2>

<h2 align="center">
    <code> go --v </code>
</h2>

<br>
<h3 align="center">(project only supports post/get requests for now, but will extend later on).</h3>
<br>
<h1 align="center">get started:</h2>

<h2 align="center">setting headers while testing:</h2>

<h3 align="center"> 
    <code> go run . [-H key:value] [URL] </code>
</h>


<br>
<h2 align="center">testing post/get requests:</h2>
<h3 align="center">
    <code> go run . [-x POST/GET] [URL] </code>
</h3>

<br>

<h2 align="center">with streaming enabled: (streaming gives live response data back.)</h2>

<h3 align="center">
    <code> go run . [-H key:value] [-stream true/false] [URL] </code>
</h3>

<h3 align="center">(important to note that you can include 1 of many flags while testing aswell, but make sure to always include the URL at the end).</h3>

<br>

<h2 align="center">include response body logging:</h2>

<h3 align="center">
    <code> go run . [-b true/false -s 1024] [URL] </code>
</h3>

<h4 align="center">(important to note you can set a limit on the response body size to log (positive values only), as with the -s flag, if you want it to be the default (1024), simply do not include the flag).</h4>

<br>

<h2 align="center">session mode:</h2>

<h3 align="center">Instead of manually pasting, you can enter session mode to dynamically store urls, headers, etc in a variable and then use that variable to test API's (in which the variable will hold the header/url etc).</h3>

<br>
<h3 align="center">to enter session mode, run the following:</h3>

<h3 align="center">
    <code> go run . [-session true/false] </code>
</h3>

<h3 align="center">when entering session mode, the url does not need to be present at the end.</h3>
<br>
<h2 align="center">session mode tutorial:</h2>
<h3 align="center">start off by setting a variable to hold a url/header or such.</h3>
<h3 align="center">
    <code> VAR [VarName] [Value]</code>
</h3>
<h3 align="center">and to retrieve a variable:</h3>
<h3 align="center">
    <code> GET [VarName] </code>
</h3>
<h3 align="center">and it returns the value of whatever the variable stores.</h3>
<h3 align="center">to delete a variable, run:</h3>
<h3 align="center">
    <code> DEL [VarName] </code>
</h3>
<h3 align="center">to see all stored variables:</h3>
<h3 align="center">
    <code> FETCH </code>
</h3>
<h3 align="center">view all available commands whilst your in session mode:</h3>
<h3 align="center">
    <code> HELP </code>
</h3>
<br>
<h3 align="center">using the variables you stored for API testing:</h3>

<h4 align="center">(for standard get requests):</h4>
<h3 align="center">
    <code> TEST [VarName] </code>
</h3>

<h4 align="center">if you a custom header aswell, do so in the following format:</h4>
<h3 align="center">
    <code> TEST [VarName] [-h headerName:value] </code>
</h3>

<h4 align="center">you can also make POST/GET requests as so:</h4>
<h3 align="center">
    <code> TEST [VarName] [-x GET] </code>
</h3>
<h4 align="center">and post request like so:</h4>
<h3 align="center">
    <code> TEST [VarName] [-x POST -D [bodyData]] </code>
</h3>
<h4 align="center">for post requests, when setting a request body you can do -F for form related data, or -D for normal json data as shown above, you can use the -F flag like so:</h4>
<h3 align="center">
    <code> TEST [VarName] [-x POST  -F [title:value]] </code>
</h3>
<br>
<h3 align="center">and finally, to exit:</h3>
<h3 align="center">
    <code> EXIT </code>
</h3>
<br>
<h4 align="center">and make sure you use the exact format as shown, (title:value).</h4>
<h4 align="center">also important, you can include multiple flag options at once when testing (except -F and -D for -x POST).</h4>

<h4 align="center">important to note, VarName should be a variable assigned with a valid url, however for the purpose of session mode, you can store headers/body data/urls in variables to utilise them when running the TEST cmd.</h4>
