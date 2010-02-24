import err
import copy
import string
        

class Stream:
    "A parse stream"
    
    def __init__(self):
        # The set of streams.
        self.streams=[]
        
        # A "pointer" to the data in the current stream.
        self.data=None 
        
        # The index of the stream we are currently in
        self.index=0
        
        # Pointer to the index of that data that we are in.
        self.at=-1
        
        # The current row
        self.row=1
        
        # The current col
        self.col=1
        
        # The stack of contexts.
        self.contexts=[]
        
        # Set to true if we are in a context manager context
        self.in_context_mgr=False
        self.ctx_manager_rollback=False
        
        
        
    def __enter__(self):
        "Enter the context manager."
        self.begin_transaction()
        self.in_context_mgr=True
        
        # Toss out the leading whitespace.
        self.skip_leading_whitespace()
               
        return self
        
    def __exit__(self, type, value, traceback):
        "When leaving the context, check to see what to do."
        if len(self.contexts) == 0:
            self.in_context_mgr=False
        
        if self.ctx_manager_rollback:
            self._rollback()
        else:
            self._commit()
                    
        self.ctx_manager_rollback=False
            

    def merge(self, data, filename=None):
        "Merges the given data into the data stream at our current point in the stream."
        
        if data==None:
            return
        
        if self.data == None:
            # Generate a new stream tag        
            d = {"filename" : filename if filename!=None else d1["filename"],
                 "data" : data }
            
            self.streams.append(d)
            self.data=data
        else:            
            # Split the current stream
            d1=None
            partition = self.at+1
            if partition < len(self.data):
                tmp=self.data[0:partition] 
                tmp2=self.data[partition:] 
            
                # Duplicate it
                d1=copy.copy(self.streams[self.index])
                d1["data"] = tmp2
                self.streams[self.index]["data"] = tmp
            
            # Generate a new stream tag        
            d2= {"filename" : filename if filename!=None else self.streams[self.index]["filename"],
                 "data" : data }
            
            # Insert the new tags in reverse order so that the new data is before the split point
            if d1!=None: self.streams.insert(self.index+1, d1)
            self.streams.insert(self.index+1, d2)
        
    def begin_transaction(self):
        "Save a context onto the stack"        
        self.contexts.append((self.streams, self.data, self.getMarker()))
        
    def _commit(self):
        "Performs the work of committing the stream state."
        self.contexts.pop()
        
    def _rollback(self):
        "Performs the work of rolling back the stream state."
        ctx=self.contexts[-1]
        self.streams=ctx[0]
        self.data=ctx[1]
        self.setMarker(ctx[2])
        self.contexts.pop()
        
    def commit(self):
        "Pop the last save context, committing to the current state"
        if self.in_context_mgr:            
            self.ctx_manager_rollback=False
            return 
                   
        self._commit()
        
    def rollback(self):
        "Restore a saved context from the stack."
        if self.in_context_mgr:            
            self.ctx_manager_rollback=True
            return 
        
        self._rollback()                                            
    
    def getMarker(self):
        "Returns a marker to the current stream pointer."
        return { "index" : self.index, "at" : self.at, "row" : self.row, "col" : self.col}
    
    def setMarker(self, m):
        "Sets the stream position from the marker."
        self.at=m["at"]
        self.row=m["row"]
        self.col=m["col"]
        self.index=m["index"]
        
    def getLoc(self):
        "Returns a location item."
        return err.location(self.streams[self.index]["filename"], self.row, self.col)
    
    def skip_leading_whitespace(self):
        cont=True
        while cont:
            c = self.peek()
            if (c!=None) and (c in string.whitespace):
                self.read() # consume and throw away the whitespace characters
            else:
                cont=False        
        
    def peek(self):
        "Reads a single character from the stream w/o consuming it. Returns None if there is nothing left to read."
        at=self.at+1
        
        if at>=len(self.data):
           index=self.index+1
           at=0
            
           if index>=len(self.streams):
                  return None
           
        else:
           index=self.index
           
        data=self.streams[index]["data"]
            
        return data[at]
            
        
    def read(self):
        "Reads a single character from the stream. Returns None if there is nothing left to read."
        self.at+=1
        
        if self.at>=len(self.data):
           self.index+=1
           self.at=0
            
           if self.index>len(self.streams):
                  return None
                 
           self.data=self.streams[self.index]["data"]
           if self.at >= len(self.data):
                   return None
           
        if self.data[self.at]=="\n": 
            self.row+=1
            self.col=1
        else:
            self.col+=1
            
        return self.data[self.at]
    
def new():
    return Stream()
