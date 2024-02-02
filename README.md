This project uses a pre-built GNU [sed](https://www.gnu.org/software/sed/). © GNU Project, GNU license.
This project uses sed definition files from the [Mulval Project](https://github.com/fiware-cybercaptor/mulval). © Xinming Ou, GNU license. 

```shell
NAME:
   AGG - Mulval-compatible attack graph generator

USAGE:
   attack_graph_generator command [command options]

VERSION:
   0.1.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version

   GENERATION

   --cycle, -c             whether cycles are permitted (does not guarantee them) (default: false)
   --relaxed               relax the constraint so that AND node can have multiple outgoing edges (default: false)
   --seed value, -s value  random seed (default: current Unix epoch in seconds)

   GRAPH

   --edge value, -e value  number of edges (default: 0)
   --node value, -n value  number of OR, PF and AND nodes (in this order)

   OUTPUT

   --graph, -g           generate a pdf rendition of the attack graph (default: false)
   --outdir DIR, -o DIR  output to DIR (default: current directory)
```