---
author: Yulian Kuncheff
date: 2022-08-23T23:12:00Z
draft: false
slug: stable-diffusion-fedora-amd
title: Stable Diffusion on Fedora 36 and AMD 6750XT
type: blog
tags:
  - programming
  - ai
  - stable diffusion
---
Yesterday, Stability.ai finally released their Stable Diffusion model to the public. Having played with NightCafe, Midjourney, and getting access to Dall-E 2, I was excited to be able to run this locally on my own computer. I have a AMD Threadripper 2990WX with an AMD Radeon 6750XT video card. I figured this was plenty powerful enough to get local generation so I didn't have to pay others when I have plenty of horsepower at home eating up electricity as it is. So I am just going to dive in and keep this short and sweet.

## Getting all the components

Firstly we need to get all the pieces to the puzzle, there aren't many, but it does take a bit to get it all figured out and sorted. I encourage you to clone or extract all of this in a central working directory for easy usage.

1. Download Anaconda 3
   1. For Fedora 36, I downloaded it from `https://www.anaconda.com`, just download the linux installer
   2. Once downloaded, go to where it was downloaded and run the installer using the terminal.
   3. After it is installed, you should have a `~/anaconda3` directory if using the default settings.
   4. Add `~/anaconda3` to your current path or restart your terminal (it tries to add it to your .bashrc or .zshrc at the end of the run).
      1. `export PATH="~/anaconda3/bin:$PATH"`
      2. You can test if it worked by running `conda` and seeing help output instead of an error.
2. Clone the Stable Diffusion source code.
   1. It can be found at [https://github.com/CompVis/stable-diffusion](https://github.com/CompVis/stable-diffusion)
3. Go to the HuggingFace repository for the model and downlaod the model
   1. The model is located at [https://huggingface.co/CompVis/stable-diffusion](https://huggingface.co/CompVis/stable-diffusion)
   2. You will need to signup for Hugging Face so that you can accept the terms to the repo.
   3. Once you have signed up and accepted, scroll down and click the repo for [stable-diffusion-v-1-4-original](https://huggingface.co/CompVis/stable-diffusion-v-1-4-original)
      1. This link takes you directly to the v1.4 model, keep an eye on the main huggingface page for new versions and update accordingly.
   4. At the top, hit the `Files and Versions` tab, and you will see what amounts to a Git repository. You can clone it through git, but for quickness, we will just download directly.
   5. Find the file with the `LFS` label next to it, for v1.4 its called `sd-v1-4.ckpt` and click to download it. It is a large file, so make go grab a drink or a snack depending on your download speeds.

There is an Optimized version of the Stable Diffusion txt2img script, that uses only 4gb of ram versus the 8gb+ of the primary script. This script was made by a researcher that had earlier access to the models and scripts. Stability.ai has mentioned that they will be releasing optimized scripts and such that will even run on a raspberry pi, but they have not yet, nor have they released official AMD instructions or support. THe following information is extra if you need this. If you have an Nvidia card, you can skip the AMD stuff, and if you have a 3090, you can even skip the optimized script, but I recommend the optimized version because even though it takes a bit longer (15 seconds versus 7-9 seconds for 6 images), its easier gets us the default 512x512 size or even higher, and works just as well.

1. Download the optimized version from the link below.
   1. For posterity, I did not write any of this, I watched a video on how to set this up on Windows, where a Google Drive with the scripts was linked. The zip had batch files Windows, and used custom Windows paths for the model and inference yaml. I changed this for Linux compatibility
      1. Video: [https://www.youtube.com/watch?v=z99WBrs1D3g](https://www.youtube.com/watch?v=z99WBrs1D3g)
      2. Google Drive: [https://drive.google.com/file/d/1z7lDaItYAm-3zNSTCIm1nRVkn8As-Wh3/view](https://drive.google.com/file/d/1z7lDaItYAm-3zNSTCIm1nRVkn8As-Wh3/view)
   2. You can get my optimized version from my github repo with the code changes.
      1. Github repository: [https://github.com/Daegalus/stable-diffusion-optimized](https://github.com/Daegalus/stable-diffusion-optimized)

## Setting things up the files and structure

For the purposes of examples, we will assume you cloned and downloaded everything into `$SD_WORKSPACE`, which can be any folder where you have everything located.

From here we want to navigate to the stable diffusion cloned repository.

```sh
cd $SD_WORKSPACE/stable-diffusion
```

Now we want to copy the model into the correct folder. This assumes you downloaded the `ckpt` checkpoint file into `$SD_WORKSPACE`, if you cloned it with git, it will probably be under its own subfolder.

```sh
mkdir -p models/ldm/stable-diffusion-v1
cp ../sd-v1-4.ckpt models/ldm/stable-diffusion-v1
```

Now, if you have Nvidia and you chose to not use the optimized third party scripts, you are done and can move onto the next section `Running Stable Diffusion`.

### Optimized Stable Diffusion

While still in the stable-diffusion folder, we can copy the new scripts to the folder.

```sh
cp ../stable-diffusion-optimized/* .
```

All set here.

### Adding AMD ROCm packages

Ok, this was the trickier part, and i had a lot of trial and error but I finally got it all working and shouln't be too hard.

I am using Fedora 36 (read footnotes for why this over other distributions). So the commands below or for this distro, please look up the appropriate packages for your distribution.

Since I am not super in the AI dev world, I wasn't sure which components we needed, so I installed all the ROCm packages available. If you find better info, adjust as needed. AMD has official repos on their website.

```sh
sudo dnf install rocm-clinfo rocm-comgr rocm-device-libs rocm-opencl rocm-runtime rocm-smi rocminfo
```

This will install ROCm 5.2.1 as of this writing.

## Setting up the rest of the environment

Now it is almost time to get it all hooked up. First thing we need to do is edit the `environment.yaml` file

1. Change `pytorch` version to `1.12.1`
2. Change `pytorch-lightning` version to `1.5.2`
3. Save and close.

Now we just run the following command to create our environment

```sh
conda env create -f environment.yaml
```

Wait for it to finish. And then finally we need to activate the environment. In general, make sure this is activated when running further commands, as they all need to run in the environment.

```sh
conda activate ldm
```

You can stop here if you are using Nvidia. You can now run `low_ram_nvidia.sh` if using the optimized scripts, or `scripts/txt2img.py` if not.

## Installing AMD PyTorch

We are going to install PyTorch with ROCm 5.1.1 support, which will work fine on the ROCm 5.2.1 we have installed.

Go to [https://pytorch.org/get-started/locally/](https://pytorch.org/get-started/locally/) and select the following comibnation: Stable, Linux, Pip, Python, ROCm 5.1.1. (ROCm can be a different version if reading in the future.)

**IMPORTANT**: Make sure you are inside your `ldm` conda environemnt with `conda activate ldm` before running the pytorch install below.

This will give you a command to run, as of this writing, the command is as follows:

```sh
// DONT RUN THIS YET, READ ON.
pip3 install torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/rocm5.1.1
```

But for Fedora, `pip3` will give us the wrong results, make sure you change the command to use `pip` as follows:

```sh
pip install torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/rocm5.1.1
```

Wait for it to finish, and we should be almost set, lets verify our setup.

Run the following commands. the `conda activate ldm` will create our conda environment ENV variables, so the right python installation inside conda is used.

```sh
conda activate ldm
python -c "import torch; print(torch.cuda.is_available())"
```

This should return `True` if everything is working. Otherwise it will show an error or `False`.

If you get an error like below, and you are running an AMD graphics card, primarily a 6700XT or 6750XT, we will need to add an env variable for the above test to work.

```sh
"hipErrorNoBinaryForGpu: Unable to find code object for all current devices!"
```

If we do an `export HSA_OVERRIDE_GFX_VERSION="10.3.0"`, and rerun the python command above, it should work.
The reasoning behind this is that most consumer cards aren't officially supported, and the only one that is included in support for ROCm is for the `gfx1030`. The 6700 and 6750 are `gfx1031`. While compatible, not directly supported, so we force it.
To figure out which `gfx` you have, you can run `rocminfo | grep gfx`.

This should get you to the point where you can run Stable diffusion successfully.

## Running Stable Diffusion

Moment of truth. Run the following commands depending on your setup:

1. AMD, optimized scripts: `low_ram.sh --prompt "A happy dog"`
   1. If you see an error about a `gfx1030_20.kdb`, you can ignore it for now. Fedora doesn't have the `miopen` package with them, and the official `rhel9` repo, wouldn't load for me. This is a warning, it can be ignored.
2. AMD, unoptimized scripts: `scripts/text2img.py --prompt "A happy dog"`
   1. You might need to lower the dimensions if you get memory errors. Add `--W 384 --H 384` to your command.
3. Nvidia, no optimized scripts: `scripts/txt2img.py --prompt "A happy dog"`
4. Nvidia, with optimized scripts: `low_ram_nvidia.sh --prompt "A happy dog"`

Happy generating!

If you find any errors in my commands above, feel free to reach out to me. Contact info is at the top right of the blog.
