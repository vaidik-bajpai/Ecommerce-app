/*
  Warnings:

  - You are about to drop the column `cartId` on the `Product` table. All the data in the column will be lost.
  - Added the required column `productIds` to the `Cart` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE "Product" DROP CONSTRAINT "Product_cartId_fkey";

-- DropIndex
DROP INDEX "Cart_id_userId_key";

-- AlterTable
ALTER TABLE "Cart" ADD COLUMN     "productIds" INTEGER NOT NULL,
ALTER COLUMN "id" DROP DEFAULT;
DROP SEQUENCE "Cart_id_seq";

-- AlterTable
ALTER TABLE "Product" DROP COLUMN "cartId";

-- AddForeignKey
ALTER TABLE "Cart" ADD CONSTRAINT "Cart_productIds_fkey" FOREIGN KEY ("productIds") REFERENCES "Product"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
