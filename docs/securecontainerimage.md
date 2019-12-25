# Overview
SecureContainerImage is the secure VM (SVM) image consisting of kernel and initrd.
This is bundled in a container image which is then used by the Raksh Operator to
provision secure container

Following are the steps used for creation of SecureContainerImage

1. Create initrd
   1. Initial rootfs + agent + skopeo  
   2. With PEF  
      * Add lockboxes (secrets and keys encrypted with public keys of the PEF system)  
      * Add nonce  
      * Encrypt rootfs  
   3. Without PEF  (development)  
      * Key stored in initrd OR  
      * Key stored in Vault  
      * Add nonce  
2. Create container image consisting of initrd and guest kernel
