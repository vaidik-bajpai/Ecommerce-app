/*
  Warnings:

  - You are about to drop the column `productIds` on the `Cart` table. All the data in the column will be lost.

*/
-- DropForeignKey
ALTER TABLE "Cart" DROP CONSTRAINT "Cart_productIds_fkey";

-- AlterTable
ALTER TABLE "Cart" DROP COLUMN "productIds";

-- AlterTable
ALTER TABLE "Product" ADD COLUMN     "cartId" INTEGER;

-- AddForeignKey
ALTER TABLE "Product" ADD CONSTRAINT "Product_cartId_fkey" FOREIGN KEY ("cartId") REFERENCES "Cart"("id") ON DELETE SET NULL ON UPDATE CASCADE;
